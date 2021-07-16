package ethereum

import (
	"context"
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func onNewBlock(ctx context.Context, latestBlockHeight uint64, client *ethclient.Client, kafkaProducer *sarama.SyncProducer) error {
	return db.IterateEthTx(func(key []byte, value []byte) error {
		var ethTx db.EthereumBroadcastedWrapTokenTransaction
		err := json.Unmarshal(value, &ethTx)
		if err != nil {
			logging.Fatal("Failed to unmarshal EthTransaction [ETH onNewBlock]:", err)
			return err
		}
		txReceipt, err := client.TransactionReceipt(ctx, ethTx.TxHash)
		if err != nil {
			if txReceipt == nil && err == ethereum.NotFound {
				logging.Info("Pending broadcast eth tx:", ethTx.TxHash)
			} else {
				logging.Error("Receipt fetch failed [onNewBlock] eth tx:", ethTx.TxHash.String(), "Error:", err)
			}
		} else {
			deleteTx := false
			if txReceipt.Status == 0 {
				logging.Warn("Broadcast eth tx failed, Hash:", ethTx.TxHash.String(), "Block:", txReceipt.BlockNumber.Uint64())
				for _, msg := range ethTx.Messages {
					msgBytes, err := json.Marshal(msg)
					if err != nil {
						return err
					}
					err = utils.ProducerDeliverMessage(msgBytes, utils.ToEth, *kafkaProducer)
					if err != nil {
						logging.Error("Failed to add msg to kafka queue [ETH onNewBlock] ToEth, Message:", msg, "Error:", err)
						return err
					}
				}
				deleteTx = true
			} else {
				confirmedBlocks := latestBlockHeight - txReceipt.BlockNumber.Uint64()
				if confirmedBlocks >= 12 {
					logging.Info("Broadcast eth tx successful. Hash:", ethTx.TxHash, "Block:", txReceipt.BlockNumber.Uint64(), "Confirmed blocks:", confirmedBlocks)
					deleteTx = true
				} else {
					logging.Info("Broadcast eth tx confirmation pending. Hash:", ethTx.TxHash, "Block:", txReceipt.BlockNumber.Uint64(), "Confirmed blocks:", confirmedBlocks)
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
