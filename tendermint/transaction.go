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

func handleTxSearchResult(clientCtx client.Context, resultTxs []*tmCoreTypes.ResultTx, kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec) error {
	for _, transaction := range resultTxs {
		err := collectAllWrapAndRevertTxs(clientCtx, transaction)
		if err != nil {
			logging.Error("Failed to process tendermint transaction:", transaction.Hash.String())
			return err
		}
	}
	wrapOrRevert(kafkaProducer, protoCodec)
	return nil
}

func collectAllWrapAndRevertTxs(clientCtx client.Context, txQueryResult *tmCoreTypes.ResultTx) error {
	if txQueryResult.TxResult.GetCode() == 0 {
		// Should be used if txQueryResult.Tx is string
		//decodedTx, err := base64.StdEncoding.DecodeString(txQueryResult.Tx)
		//if err != nil {
		//	return txMsgs, err
		//}

		txInterface, err := clientCtx.TxConfig.TxDecoder()(txQueryResult.Tx)
		if err != nil {
			return err
		}

		transaction, ok := txInterface.(signing.Tx)
		if !ok {
			return err
		}

		memo := strings.TrimSpace(transaction.GetMemo())

		for i, msg := range transaction.GetMsgs() {
			switch txMsg := msg.(type) {
			case *banktypes.MsgSend:
				if txMsg.ToAddress == configuration.GetAppConfig().Tendermint.GetPStakeAddress() {
					exists := db.CheckTendermintIncomingTxProduced(txQueryResult.Hash, uint(i))
					if !exists {
						err = db.SetTendermintIncomingTx(db.TendermintIncomingTx{
							TxHash:          txQueryResult.Hash,
							MsgIndex:        uint(i),
							ProducedToKafka: false,
							Msg:             *txMsg,
							Memo:            memo,
						})
						if err != nil {
							return err
						}
					}
				}
			default:
			}
		}
	}
	return nil
}

func wrapOrRevert(kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec) {
	txs, err := db.GetProduceToKafkaTendermintTxs()
	if err != nil {
		logging.Fatal(err)
	}
	for _, tx := range txs {
		validEthMemo := goEthCommon.IsHexAddress(tx.Memo)
		var ethAddress goEthCommon.Address
		if validEthMemo {
			ethAddress = goEthCommon.HexToAddress(tx.Memo)
		}

		if tx.Memo != "DO_NOT_REVERT" {
			amount := sdk.ZeroInt()
			refundCoins := sdk.NewCoins()
			for _, coin := range tx.Msg.Amount {
				if coin.Denom == configuration.GetAppConfig().Tendermint.PStakeDenom {
					amount = coin.Amount
				} else {
					refundCoins = refundCoins.Add(coin)
				}
			}
			fromAddress, err := sdk.AccAddressFromBech32(tx.Msg.FromAddress)
			if err != nil {
				logging.Fatal(err)
			}
			totalWrappedAmount, err := db.GetTotalTokensWrapped()
			if err != nil {
				logging.Fatal(err)
			}
			if totalWrappedAmount.GTE(getMaxLimit()) {
				logging.Info("Reverting Tendermint Tx [MAX Amount Reached]:", tx.TxHash, "MsgBytes Index:", tx.MsgIndex)
				refundCoins = refundCoins.Add(sdk.NewCoin(configuration.GetAppConfig().Tendermint.PStakeDenom, amount))
				revertCoins(tx.Msg.FromAddress, refundCoins, kafkaProducer, protoCodec)
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
				logging.Info("Received Tendermint Wrap Tx:", tx.TxHash, "MsgBytes Index:", tx.MsgIndex)
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
				tx.ProducedToKafka = true
				err = db.SetTendermintIncomingTx(tx)
				if err != nil {
					logging.Fatal(err)
				}
				accountLimiter.Amount = accountLimiter.Amount.Add(sendAmt)
				err = db.SetAccountLimiter(accountLimiter)
				if err != nil {
					logging.Fatal(err)
				}
			}
			if len(refundCoins) > 0 {
				logging.Info("Reverting left over coins: TxHash:", tx.TxHash, "MsgBytes Index:", tx.MsgIndex)
				revertCoins(tx.Msg.FromAddress, refundCoins, kafkaProducer, protoCodec)
			}
		} else {
			logging.Info("Deposited to wrap address, TxHash:", tx.TxHash, "amount:", tx.Msg.Amount.String())
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
