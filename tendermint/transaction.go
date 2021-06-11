package tendermint

import (
	"encoding/json"
	"github.com/persistenceOne/persistenceBridge/application"
	ethereum2 "github.com/persistenceOne/persistenceBridge/ethereum"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"log"
	"math/big"
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

func handleTxSearchResult(clientCtx client.Context, txSearchResult *tmCoreTypes.ResultTxSearch, kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec) error {
	for _, transaction := range txSearchResult.Txs {
		err := processTx(clientCtx, transaction, kafkaState, protoCodec)
		if err != nil {
			log.Printf("Failed to process tendermint transaction: %s\n", transaction.Hash.String())
			return err
		}
	}
	return nil
}

func processTx(clientCtx client.Context, txQueryResult *tmCoreTypes.ResultTx, kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec) error {
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
				var amount *big.Int
				for _, coin := range txMsg.Amount {
					if coin.Denom == application.GetAppConfiguration().PStakeDenom {
						amount = coin.Amount.BigInt()
						break
					}
				}
				if txMsg.ToAddress == application.GetAppConfiguration().PStakeAddress.String() && amount != nil && validMemo {
					log.Printf("TM Tx: %s, Msg Index: %d\n", txQueryResult.Hash.String(), i)
					ethTxMsg := ethereum2.EthTxMsg{
						Address: ethAddress,
						Amount:  amount,
					}
					msgBytes, err := json.Marshal(ethTxMsg)
					if err != nil {
						panic(err)
					}
					err = utils.ProducerDeliverMessage(msgBytes, utils.ToEth, kafkaState.Producer)
					if err != nil {
						log.Printf("Failed to add msg to kafka queue: %s\n", err.Error())
					}
					log.Printf("Produced to kafka: %v, for topic %v \n", msg.String(), utils.ToEth)
				} else {
					msg := &banktypes.MsgSend{
						FromAddress: txMsg.ToAddress,
						ToAddress:   txMsg.FromAddress,
						Amount:      txMsg.Amount,
					}
					msgBytes, err := protoCodec.MarshalInterface(sdk.Msg(msg))
					if err != nil {
						panic(err)
					}
					err = utils.ProducerDeliverMessage(msgBytes, utils.ToTendermint, kafkaState.Producer)
					if err != nil {
						log.Printf("Failed to add msg to kafka queue: %s\n", err.Error())
					}
					log.Printf("Produced to kafka: %v, for topic %v\n", msg.String(), utils.ToTendermint)
				}
			default:

			}
		}
	}

	return nil
}