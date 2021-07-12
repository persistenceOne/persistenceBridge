package tendermint

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/outgoingTx"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"log"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	goEthCommon "github.com/ethereum/go-ethereum/common"
	tmCoreTypes "github.com/tendermint/tendermint/rpc/core/types"
)

//func handleTxEvent(clientCtx client.Context, txEvent tmTypes.EventDataTx, kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec) {
//	if txEvent.Result.Code == 0 {
//		_ = processTx(clientCtx, txEvent.Tx, kafkaState, protoCodec)
//	}
//}

func handleTxSearchResult(clientCtx client.Context, txSearchResult *tmCoreTypes.ResultTxSearch, kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec) error {
	for _, transaction := range txSearchResult.Txs {
		err := processTx(clientCtx, transaction, kafkaProducer, protoCodec)
		if err != nil {
			log.Printf("Failed to process tendermint transaction: %s\n", transaction.Hash.String())
			return err
		}
	}
	return nil
}

func processTx(clientCtx client.Context, txQueryResult *tmCoreTypes.ResultTx, kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec) error {
	if txQueryResult.TxResult.GetCode() == 0 {
		// Should be used if txQueryResult.Tx is string
		//decodedTx, err := base64.StdEncoding.DecodeString(txQueryResult.Tx)
		//if err != nil {
		//	log.Fatalln(err.Error())
		//}

		txInterface, err := clientCtx.TxConfig.TxDecoder()(txQueryResult.Tx)
		if err != nil {
			log.Fatalln(err.Error())
		}

		transaction, ok := txInterface.(signing.Tx)
		if !ok {
			log.Fatalln("Unable to parse transaction into signing.Tx")
		}

		memo := strings.TrimSpace(transaction.GetMemo())
		validMemo := goEthCommon.IsHexAddress(memo)
		var ethAddress goEthCommon.Address
		if validMemo {
			ethAddress = goEthCommon.HexToAddress(memo)
		}

		for i, msg := range transaction.GetMsgs() {
			switch txMsg := msg.(type) {
			case *banktypes.MsgSend:
				if txMsg.ToAddress == configuration.GetAppConfig().Tendermint.GetPStakeAddress() {
					if memo != "DO_NOT_REVERT" {
						amount := sdk.ZeroInt()
						refundCoins := sdk.NewCoins()
						for _, coin := range txMsg.Amount {
							if coin.Denom == configuration.GetAppConfig().Tendermint.PStakeDenom {
								amount = coin.Amount
							} else {
								refundCoins = refundCoins.Add(coin)
							}
						}
						fromAddress, err := sdk.AccAddressFromBech32(txMsg.FromAddress)
						if err != nil {
							log.Fatalln(err)
						}
						accountLimiter, totalAccounts := db.GetAccountLimiterAndTotal(fromAddress)
						if totalAccounts >= getMaxLimit() {
							log.Printf("REVERT TM %s, Msg Index %d. MAX Account Limit Reached\n", txQueryResult.Hash.String(), i)
							revertCoins(txMsg.FromAddress, txMsg.Amount, kafkaProducer, protoCodec)
							continue
						}
						sendAmt, refundAmt := beta(accountLimiter, amount)
						if refundAmt.GT(sdk.ZeroInt()) {
							refundCoins = refundCoins.Add(sdk.NewCoin(configuration.GetAppConfig().Tendermint.PStakeDenom, refundAmt))
						}
						if sendAmt.GT(sdk.ZeroInt()) && validMemo {
							log.Printf("RECEIVED TM WRAP TX: %s, Msg Index: %d\n", txQueryResult.Hash.String(), i)
							ethTxMsg := outgoingTx.WrapTokenMsg{
								Address: ethAddress,
								Amount:  sendAmt.BigInt(),
							}
							msgBytes, err := json.Marshal(ethTxMsg)
							if err != nil {
								panic(err)
							}
							log.Printf("Adding wrap token msg to kafka producer ToEth, from: %s, to: %s, amount: %s\n", fromAddress.String(), ethAddress.String(), sendAmt.String())
							err = utils.ProducerDeliverMessage(msgBytes, utils.ToEth, *kafkaProducer)
							if err != nil {
								log.Fatalf("Failed to add msg to kafka queue ToEth: %s [TM Listener]\n", err.Error())
							}
							accountLimiter.Amount = accountLimiter.Amount.Add(sendAmt)
							err = db.SetAccountLimiter(accountLimiter)
							if err != nil {
								panic(err)
							}
						}
						if len(refundCoins) > 0 {
							log.Printf("REVERT LEFT OVER COINS: TM %s Msg Index: %d.\n", txQueryResult.Hash.String(), i)
							revertCoins(txMsg.FromAddress, refundCoins, kafkaProducer, protoCodec)
						}
					} else {
						log.Printf("TM Tx %s deposited %s to wrap address", txQueryResult.Hash.String(), txMsg.Amount.String())
					}
				}
			default:

			}
		}
	}

	return nil
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

func getMaxLimit() int {
	currentTime := time.Now()
	// 19th July, 2021
	if currentTime.Unix() < 1626696000 {
		return 50000
	}
	// 26th July, 2021
	if currentTime.Unix() < 1627300800 {
		return 65000
	}
	// 2nd August, 2021
	if currentTime.Unix() < 1627905600 {
		return 80000
	}
	// 9th August, 2021
	if currentTime.Unix() < 1628510400 {
		return 95000
	}
	// 16th August, 2021
	if currentTime.Unix() < 1629115200 {
		return 110000
	}
	// 23rd August, 2021
	if currentTime.Unix() < 1629720000 {
		return 125000
	}
	return 2147483646
}

func revertCoins(toAddress string, coins sdk.Coins, kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec) {
	msg := &banktypes.MsgSend{
		FromAddress: configuration.GetAppConfig().Tendermint.GetPStakeAddress(),
		ToAddress:   toAddress,
		Amount:      coins,
	}
	msgBytes, err := protoCodec.MarshalInterface(msg)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("REVERT: adding send coin msg to kafka producer MsgSend, to: %s, amount: %s\n", toAddress, coins.String())
	err = utils.ProducerDeliverMessage(msgBytes, utils.MsgSend, *kafkaProducer)
	if err != nil {
		log.Fatalf("Failed to add msg to kafka queue MsgSend: %s [TM Listener (REVERT)]\n", err.Error())
	}
}
