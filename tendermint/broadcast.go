package tendermint

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/relayer/relayer"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"log"
)

func onNewBlock(ctx context.Context, chain *relayer.Chain, kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec) error {
	return db.IterateTmTx(func(key []byte, value []byte) error {
		fmt.Println("GOT KEY: " + string(key))
		var tmTx db.TMTransaction
		err := json.Unmarshal(value, &tmTx)
		if err != nil {
			log.Fatalln("Failed to unmarshal TMTransaction : ", err)
		}
		txSearchResult, err := chain.Client.TxSearch(ctx, fmt.Sprintf("tx.hash='%s'", tmTx.TxHash), true, nil, nil, "asc")
		if err != nil {
			if txSearchResult != nil && len(txSearchResult.Txs) == 1 {
				if txSearchResult.Txs[0].TxResult.GetCode() != 0 {
					log.Printf("Tx %s to TM failed, code: %d, log: %s\n", tmTx.TxHash, txSearchResult.Txs[0].TxResult.GetCode(), txSearchResult.Txs[0].TxResult.Log)
					for _, msg := range tmTx.Messages {
						msgBytes, err := protoCodec.MarshalInterface(msg)
						if err != nil {
							log.Fatalln("Failed to generate msgBytes: ", err)
						}
						log.Printf("Adding failed tx %s to kafka producer %s: %s\n", tmTx.TxHash, utils.ToTendermint, msg.String())
						err = utils.ProducerDeliverMessage(msgBytes, utils.ToTendermint, kafkaState.Producer)
						if err != nil {
							log.Fatalf("Failed to add msg to kafka topic %s queue: %s\n", utils.ToTendermint, err.Error())
						}

					}
					return db.Delete(key)
				}
			} else {
				log.Fatalf("unknown txSearchResult: %v\n", txSearchResult)
			}
		} else {
			log.Printf("tx hash search failed: %s\n", err)
		}

		return nil
	})
}
