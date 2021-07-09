package tendermint

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/outgoingTx"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"log"
	"strings"

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
				if txMsg.ToAddress == configuration.GetAppConfig().Tendermint.PStakeAddress && memo != "DO_NOT_REFUND" {
					amount := sdk.ZeroInt()
					var refundCoins sdk.Coins
					for _, coin := range txMsg.Amount {
						if coin.Denom == configuration.GetAppConfig().Tendermint.PStakeDenom {
							amount = coin.Amount
							break
						} else {
							refundCoins = append(refundCoins, coin)
						}
					}
					fromAddress, err := sdk.AccAddressFromBech32(txMsg.FromAddress)
					if err != nil {
						log.Fatalln(err)
					}
					exists, accountLimiter, totalAccounts := db.CheckExistsAndGetTotalAccounts(fromAddress)
					if totalAccounts >= getMaxLimit() {
						refundAmount(txMsg.FromAddress, txMsg.Amount, kafkaProducer, protoCodec)
						continue
					}
					sendAmt, refundAmt := beta(exists, accountLimiter, amount)
					if refundAmt.GT(sdk.ZeroInt()) {
						refundCoins = append(refundCoins, sdk.NewCoin(configuration.GetAppConfig().Tendermint.PStakeDenom, refundAmt))
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
						err = utils.ProducerDeliverMessage(msgBytes, utils.ToEth, *kafkaProducer)
						if err != nil {
							log.Printf("Failed to add msg to kafka queue: %s\n", err.Error())
						}
						log.Printf("Produced to kafka: %v, for topic %v \n", msg.String(), utils.ToEth)
						accountLimiter.Amount = accountLimiter.Amount.Add(sendAmt)
						err = db.SetAccountLimiter(accountLimiter)
						if err != nil {
							panic(err)
						}
					}
					if len(refundCoins) > 0 {
						refundAmount(txMsg.FromAddress, refundCoins, kafkaProducer, protoCodec)
					}
				}
			default:

			}
		}
	}

	return nil
}

func beta(exists bool, limiter db.AccountLimiter, amount sdk.Int) (sdk.Int, sdk.Int) {
	maxAmt := sdk.NewInt(500000000)
	sendAmount := amount
	refundAmt := sdk.ZeroInt()
	if exists {
		sent := limiter.Amount
		if sent.Add(sendAmount).GTE(maxAmt) {
			sendAmount = maxAmt.Sub(sent)
			refundAmt = amount.Sub(sendAmount)
		}
	} else {
		if amount.GTE(constants.MinimumAmount) {
			if sendAmount.GTE(maxAmt) {
				sendAmount = maxAmt
				refundAmt = amount.Sub(maxAmt)
			}
		} else {
			sendAmount = sdk.ZeroInt()
			refundAmt = amount
		}
	}
	return sendAmount, refundAmt
}

func getMaxLimit() int {
	return 10000
}

func refundAmount(toAddress string, coins sdk.Coins, kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec) {
	msg := &banktypes.MsgSend{
		FromAddress: configuration.GetAppConfig().Tendermint.PStakeAddress,
		ToAddress:   toAddress,
		Amount:      coins,
	}
	msgBytes, err := protoCodec.MarshalInterface(sdk.Msg(msg))
	if err != nil {
		log.Fatalln(err)
	}
	err = utils.ProducerDeliverMessage(msgBytes, utils.MsgSend, *kafkaProducer)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Produced to kafka: %v, for topic %v\n", msg.String(), utils.MsgSend)
}
