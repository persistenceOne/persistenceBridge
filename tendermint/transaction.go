package tendermint

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	goEthCommon "github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/outgoingTx"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
	tmCoreTypes "github.com/tendermint/tendermint/rpc/core/types"
)

type tmWrapOrRevert struct {
	txHash   string
	msgIndex int
	msg      *banktypes.MsgSend
	memo     string
}

func handleTxSearchResult(clientCtx client.Context, resultTxs []*tmCoreTypes.ResultTx, kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec) error {
	var allTxsWrapOrRevert []tmWrapOrRevert
	for _, transaction := range resultTxs {
		tmWrapOrReverts, err := collectAllWrapAndRevertTxs(clientCtx, transaction)
		if err != nil {
			logging.Error("Failed to process tendermint transaction:", transaction.Hash.String())
			return err
		}
		allTxsWrapOrRevert = append(allTxsWrapOrRevert, tmWrapOrReverts...)
	}
	wrapOrRevert(allTxsWrapOrRevert, kafkaProducer, protoCodec)
	return nil
}

func collectAllWrapAndRevertTxs(clientCtx client.Context, txQueryResult *tmCoreTypes.ResultTx) ([]tmWrapOrRevert, error) {
	var tmWrapOrReverts []tmWrapOrRevert
	if txQueryResult.TxResult.GetCode() == 0 {
		// Should be used if txQueryResult.Tx is string
		//decodedTx, err := base64.StdEncoding.DecodeString(txQueryResult.Tx)
		//if err != nil {
		//	return txMsgs, err
		//}

		txInterface, err := clientCtx.TxConfig.TxDecoder()(txQueryResult.Tx)
		if err != nil {
			return tmWrapOrReverts, err
		}

		transaction, ok := txInterface.(signing.Tx)
		if !ok {
			return tmWrapOrReverts, err
		}

		memo := strings.TrimSpace(transaction.GetMemo())

		for i, msg := range transaction.GetMsgs() {
			switch txMsg := msg.(type) {
			case *banktypes.MsgSend:
				if txMsg.ToAddress == configuration.GetAppConfig().Tendermint.GetPStakeAddress() {
					t := tmWrapOrRevert{
						txHash:   txQueryResult.Hash.String(),
						msgIndex: i,
						msg:      txMsg,
						memo:     memo,
					}
					tmWrapOrReverts = append(tmWrapOrReverts, t)
				}
			default:
			}
		}
	}
	return tmWrapOrReverts, nil
}

func wrapOrRevert(tmWrapOrReverts []tmWrapOrRevert, kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec) {
	for _, wrapOrRevertMsg := range tmWrapOrReverts {
		validEthMemo := goEthCommon.IsHexAddress(wrapOrRevertMsg.memo)
		var ethAddress goEthCommon.Address
		if validEthMemo {
			ethAddress = goEthCommon.HexToAddress(wrapOrRevertMsg.memo)
		}

		if wrapOrRevertMsg.memo != "DO_NOT_REVERT" {
			amount := sdk.ZeroInt()
			refundCoins := sdk.NewCoins()
			for _, coin := range wrapOrRevertMsg.msg.Amount {
				if coin.Denom == configuration.GetAppConfig().Tendermint.PStakeDenom {
					amount = coin.Amount
				} else {
					refundCoins = refundCoins.Add(coin)
				}
			}
			fromAddress, err := sdk.AccAddressFromBech32(wrapOrRevertMsg.msg.FromAddress)
			if err != nil {
				logging.Fatal(err)
			}
			totalWrappedAmount, err := db.GetTotalTokensWrapped()
			if err != nil {
				logging.Fatal(err)
			}
			if totalWrappedAmount.GTE(getMaxLimit()) {
				logging.Info("Reverting Tendermint Tx [MAX Amount Reached]:", wrapOrRevertMsg.txHash, "Msg Index:", wrapOrRevertMsg.msgIndex)
				refundCoins = refundCoins.Add(sdk.NewCoin(configuration.GetAppConfig().Tendermint.PStakeDenom, amount))
				revertCoins(wrapOrRevertMsg.msg.FromAddress, refundCoins, kafkaProducer, protoCodec)
				refundCoins = sdk.NewCoins()
				continue
			}
			accountLimiter, err := db.GetAccountLimiter(fromAddress)
			if err != nil {
				logging.Fatal(err)
			}
			sendAmt, refundAmt := beta(accountLimiter, amount)
			if refundAmt.GT(sdk.ZeroInt()) {
				refundCoins = refundCoins.Add(sdk.NewCoin(configuration.GetAppConfig().Tendermint.PStakeDenom, refundAmt))
			}
			if sendAmt.GT(sdk.ZeroInt()) && validEthMemo {
				logging.Info("Received Tendermint Wrap Tx:", wrapOrRevertMsg.txHash, "Msg Index:", wrapOrRevertMsg.msgIndex)
				ethTxMsg := outgoingTx.WrapTokenMsg{
					Address: ethAddress,
					Amount:  sendAmt.BigInt(),
				}
				msgBytes, err := json.Marshal(ethTxMsg)
				if err != nil {
					panic(err)
				}
				logging.Info("Adding wrap token msg to kafka producer ToEth, from:", fromAddress.String(), "to:", ethAddress.String(), "amount:", sendAmt.String())
				err = utils.ProducerDeliverMessage(msgBytes, utils.ToEth, *kafkaProducer)
				if err != nil {
					logging.Fatal("Failed to add msg to kafka queue ToEth [TM Listener]:", err)
				}
				accountLimiter.Amount = accountLimiter.Amount.Add(sendAmt)
				err = db.SetAccountLimiter(accountLimiter)
				if err != nil {
					logging.Fatal(err)
				}
			}
			if len(refundCoins) > 0 {
				logging.Info("Reverting left over coins: TxHash:", wrapOrRevertMsg.txHash, "Msg Index:", wrapOrRevertMsg.msgIndex)
				revertCoins(wrapOrRevertMsg.msg.FromAddress, refundCoins, kafkaProducer, protoCodec)
			}
		} else {
			logging.Info("Deposited to wrap address, TxHash:", wrapOrRevertMsg.txHash, "amount:", wrapOrRevertMsg.msg.Amount.String())
		}
	}
}

func beta(limiter db.AccountLimiter, amount sdk.Int) (sendAmount sdk.Int, refundAmt sdk.Int) {
	if amount.LT(sdk.NewInt(configuration.GetAppConfig().Tendermint.MinimumWrapAmount)) {
		sendAmount = sdk.ZeroInt()
		refundAmt = amount
		return sendAmount, refundAmt
	}
	maxAmt := sdk.NewInt(int64(500000000))
	sendAmount = amount
	refundAmt = sdk.ZeroInt()
	sent := limiter.Amount
	if sent.Add(sendAmount).GTE(maxAmt) {
		sendAmount = maxAmt.Sub(sent)
		refundAmt = amount.Sub(sendAmount)
	}
	return sendAmount, refundAmt
}

func getMaxLimit() sdk.Int {
	currentTime := time.Now()
	// 19th July 2021
	if currentTime.Unix() < 1626696000 {
		return sdk.NewInt(50000000000)
	}
	// 26th July 2021
	if currentTime.Unix() < 1627300800 {
		return sdk.NewInt(65000000000)
	}
	// 2nd August 2021
	if currentTime.Unix() < 1627905600 {
		return sdk.NewInt(80000000000)
	}
	// 9th August 2021
	if currentTime.Unix() < 1628510400 {
		return sdk.NewInt(95000000000)
	}
	// 16th August 2021
	if currentTime.Unix() < 1629115200 {
		return sdk.NewInt(110000000000)
	}
	// 23rd August, 2021
	if currentTime.Unix() < 1629720000 {
		return sdk.NewInt(125000000000)
	}
	// 30th August 2021
	if currentTime.Unix() < 1630324800 {
		return sdk.NewInt(140000000000)
	}
	// 6th September 2021
	if currentTime.Unix() < 1630929600 {
		return sdk.NewInt(155000000000)
	}
	// 13th September 2021
	if currentTime.Unix() < 1631534400 {
		return sdk.NewInt(170000000000)
	}
	// 20th September 2021
	if currentTime.Unix() < 1632139200 {
		return sdk.NewInt(185000000000)
	}
	return sdk.NewInt(1000000000000000)
}

func revertCoins(toAddress string, coins sdk.Coins, kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec) {
	msg := &banktypes.MsgSend{
		FromAddress: configuration.GetAppConfig().Tendermint.GetPStakeAddress(),
		ToAddress:   toAddress,
		Amount:      coins,
	}
	msgBytes, err := protoCodec.MarshalInterface(msg)
	if err != nil {
		logging.Fatal(err)
	}
	logging.Info("REVERT: adding send coin msg to kafka producer MsgSend, to:", toAddress, "amount:", coins.String())
	err = utils.ProducerDeliverMessage(msgBytes, utils.MsgSend, *kafkaProducer)
	if err != nil {
		logging.Fatal("Failed to add msg to kafka queue ToEth [TM Listener REVERT]:", err)
	}
}
