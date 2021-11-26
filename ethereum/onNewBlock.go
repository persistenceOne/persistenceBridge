package ethereum

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func onNewBlock(ctx context.Context, latestBlockHeight uint64, client *ethclient.Client, kafkaProducer *sarama.SyncProducer) error {
	return db.IterateOutgoingEthTx(func(key []byte, value []byte) error {
		var ethTx db.OutgoingEthereumTransaction
		err := json.Unmarshal(value, &ethTx)
		if err != nil {
			return fmt.Errorf("failed to unmarshal OutgoingEthereumTransaction %s [ETH onNewBlock]: %s", string(key), err.Error())
		}
		fmt.Println(ethTx.TxHash)
		transaction, _, err := client.TransactionByHash(ctx, ethTx.TxHash)
		txReceipt, err := client.TransactionReceipt(ctx, ethTx.TxHash)
		fmt.Println(txReceipt.Logs)
		sender, err := client.TransactionSender(ctx, transaction, txReceipt.BlockHash,txReceipt.TransactionIndex)
		fmt.Println(sender.Hash())
		if err!= nil {
			fmt.Println(err.Error())
		}
		if err != nil {
			if txReceipt == nil && err == ethereum.NotFound {
				logging.Info("Broadcast ethereum tx pending:", ethTx.TxHash)
			} else {
				logging.Error("Receipt fetch failed [onNewBlock] eth tx:", ethTx.TxHash.String(), "Error:", err)
			}
		} else {
			deleteTx := false
			if txReceipt.Status == 0 {
				err := logging.ErrorReason(ctx, sender, transaction, txReceipt.BlockNumber, *client)
				if err!=nil {
					return err
				}
				logging.Error("Broadcast ethereum tx failed, Hash:", ethTx.TxHash.String(), "Block:", txReceipt.BlockNumber.Uint64())
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
					logging.Info("Broadcast ethereum tx successful. Hash:", ethTx.TxHash, "Block:", txReceipt.BlockNumber.Uint64(), "Confirmed blocks:", confirmedBlocks)
					deleteTx = true
				} else {
					logging.Info("Broadcast ethereum tx confirmation pending. Hash:", ethTx.TxHash, "Block:", txReceipt.BlockNumber.Uint64(), "Confirmed blocks:", confirmedBlocks)
				}
			}
			if deleteTx {
				return db.DeleteOutgoingEthereumTx(ethTx.TxHash)
			}
			return nil
		}
		return nil
	})


}
