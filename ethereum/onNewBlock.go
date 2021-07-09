package ethereum

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"log"
)

func onNewBlock(ctx context.Context, client *ethclient.Client, kafkaState utils.KafkaState) error {
	return db.IterateEthTx(func(key []byte, value []byte) error {
		var ethTx db.EthereumBroadcastedWrapTokenTransaction
		err := json.Unmarshal(value, &ethTx)
		if err != nil {
			log.Fatalln("Failed to unmarshal EthTransaction: ", err)
		}
		txReceipt, err := client.TransactionReceipt(ctx, ethTx.TxHash)
		if err != nil {
			if txReceipt == nil && err == ethereum.NotFound {
				log.Printf("ETH TX %s is in pending transactions\n", ethTx.TxHash)
			} else {
				log.Printf("ETH TX %s receipt fetch failed: %s\n", ethTx.TxHash.String(), err)
			}
		} else {
			if txReceipt.Status == 0 {
				log.Printf("Broadacasted ethereum tx failed: %s\n", ethTx.TxHash.String())
				for _, msg := range ethTx.Messages {
					msgBytes, err := json.Marshal(msg)
					if err != nil {
						log.Fatalln("Failed to generate msgBytes: ", err)
					}
					err = utils.ProducerDeliverMessage(msgBytes, utils.ToEth, kafkaState.Producer)
					if err != nil {
						log.Fatalf("Failed to add msg to kafka topic %s queue: %s\n", utils.ToEth, err.Error())
					}
				}
			} else {
				log.Printf("Broadcasted ethereum tx %s success\n", ethTx.TxHash)
			}
			return db.DeleteEthereumTx(ethTx.TxHash)
		}
		return nil
	})
}
