package ethereum

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"log"
)

func onNewBlock(ctx context.Context, latestBlockHeight uint64, client *ethclient.Client, kafkaProducer *sarama.SyncProducer) error {
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
			deleteTx := false
			if txReceipt.Status == 0 {
				log.Printf("Broadacasted ethereum tx failed: %s\n", ethTx.TxHash.String())
				for _, msg := range ethTx.Messages {
					msgBytes, err := json.Marshal(msg)
					if err != nil {
						log.Fatalln("Failed to generate msgBytes: ", err)
					}
					err = utils.ProducerDeliverMessage(msgBytes, utils.ToEth, *kafkaProducer)
					if err != nil {
						log.Fatalf("Failed to add msg to kafka topic %s queue: %s\n", utils.ToEth, err.Error())
					}
				}
				deleteTx = true
			} else {
				confirmedBlocks := latestBlockHeight - txReceipt.BlockNumber.Uint64()
				if confirmedBlocks >= 12 {
					log.Printf("Broadcasted ethereum tx %s success. Has %d confirmed blocks\n", ethTx.TxHash, confirmedBlocks)
					deleteTx = true
				} else {
					log.Printf("Broadcasted ethereum tx %s has %d block confirmations\n", ethTx.TxHash, confirmedBlocks)
				}
			}
			if deleteTx {
				return db.DeleteEthereumTx(ethTx.TxHash)
			}
			return nil
		}
		return nil
	})
}
