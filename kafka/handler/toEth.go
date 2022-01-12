/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package handler

import (
	"encoding/json"
	"time"

	"github.com/Shopify/sarama"
	"github.com/dgraph-io/badger/v3"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/outgoingtx"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func (m MsgHandler) HandleToEth(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var kafkaMsgs []sarama.ConsumerMessage

	claimMsgChan := claim.Messages()

	ticker := time.NewTicker(configuration.GetAppConfig().Kafka.ToEth.Ticker)
	defer ticker.Stop()

	var (
		kafkaMsg *sarama.ConsumerMessage
		ok       bool
	)

ConsumerLoop:
	for {
		select {
		case <-ticker.C:
			if len(kafkaMsgs) >= configuration.GetAppConfig().Kafka.ToEth.MinBatchSize {
				break ConsumerLoop
			} else {
				return nil
			}
		case kafkaMsg, ok = <-claimMsgChan:
			if ok {
				kafkaMsgs = append(kafkaMsgs, *kafkaMsg)
				if len(kafkaMsgs) == configuration.GetAppConfig().Kafka.ToEth.MaxBatchSize {
					break ConsumerLoop
				}
			} else {
				break ConsumerLoop
			}
		}
	}

	if len(kafkaMsgs) == 0 {
		return nil
	}

	if kafkaMsg == nil {
		return ErrKafkaNilMessage
	}

	// 1.add to database
	msgBytes := make([][]byte, len(kafkaMsgs))

	for i := range kafkaMsgs {
		msgBytes[i] = kafkaMsgs[i].Value
	}

	index, err := db.AddKafkaEthereumConsume(m.DB, kafkaMsg.Offset, msgBytes)
	if err != nil {
		return err
	}

	// 2.set kafka offset so all next steps are independent of kafka consumer queue.
	session.MarkMessage(kafkaMsg, "")

	return SendBatchToEth(index, m)
}

func convertMsgBytesToEthMsg(msgBytes [][]byte) ([]db.WrapTokenMsg, error) {
	msgs := make([]db.WrapTokenMsg, len(msgBytes))

	var err error

	for i := range msgBytes {
		var msg db.WrapTokenMsg

		err = json.Unmarshal(msgBytes[i], &msg)
		if err != nil {
			return nil, err
		}

		msgs[i] = msg
	}

	return msgs, nil
}

// SendBatchToEth : Handling of msgSend
func SendBatchToEth(index uint64, handler MsgHandler) error {
	kafkaConsume, err := db.GetKafkaEthereumConsume(handler.DB, index)
	if err != nil {
		logging.Fatal(err)
	}
	msgs, err := convertMsgBytesToEthMsg(kafkaConsume.MsgBytes)
	if err != nil {
		logging.Fatal(err)
	}

	logging.Info("batched messages to send to ETH:", msgs)

	hash, err := outgoingtx.EthereumWrapAndStakeToken(handler.EthClient, msgs)
	if err != nil {
		logging.Error("Unable to do ethereum tx (adding messages again to kafka), messages:", msgs, "error:", err)

		config := utils.SaramaConfig()
		producer := utils.NewProducer(configuration.GetAppConfig().Kafka.Brokers, config)

		defer func() {
			innerErr := producer.Close()
			if innerErr != nil {
				logging.Error("failed to close producer in topic: SendBatchToEth, err:", innerErr)
			}
		}()

		err = db.DeleteKafkaEthereumConsume(handler.DB, index)
		if err != nil {
			logging.Error("Failed to delete Ethereum msg at index: ", index, " Error: ", err)
		}

		var msgBytes []byte

		for i := range msgs {
			msgBytes, err = json.Marshal(msgs[i])
			if err != nil {
				logging.Error("Failed to Marshal an unmarshalled msg")
			}

			err = utils.ProducerDeliverMessage(msgBytes, utils.ToEth, producer)
			if err != nil {
				logging.Error("Failed to add msg to kafka ToEth queue (need to do manually), message:", msgs[i], "error:", err)
			}
		}

		return err
	}

	err = db.UpdateKafkaEthereumConsumeTxHash(handler.DB, index, hash)
	if err != nil {
		logging.Fatal(err)
	}

	err = db.SetOutgoingEthereumTx(handler.DB, db.NewOutgoingETHTransaction(hash, msgs))
	if err != nil {
		logging.Fatal(err)
	}

	logging.Info("Broadcast Eth Tx hash:", hash.String())

	checkKafkaEthereumConsumeDBAndAddToRetry(handler.DB)

	return nil
}

func checkKafkaEthereumConsumeDBAndAddToRetry(database *badger.DB) {
	// all logging, no return
	kafkaEthereumConsumes, err := db.GetEmptyTxHashesETH(database)
	if err != nil {
		logging.Error(err)
	}

	if len(kafkaEthereumConsumes) > 0 {
		config := utils.SaramaConfig()
		producer := utils.NewProducer(configuration.GetAppConfig().Kafka.Brokers, config)

		defer func() {
			err := producer.Close()
			if err != nil {
				logging.Error("failed to close producer in topic: SendBatchToEth, err:", err)
			}
		}()

		for _, kafkaEthereumConsume := range kafkaEthereumConsumes {

			err = db.DeleteKafkaEthereumConsume(database, kafkaEthereumConsume.Index)
			if err != nil {
				logging.Error("Failed to delete Ethereum msg at index: ", kafkaEthereumConsume.Index, " Error: ", err)
			}

			for _, msgByte := range kafkaEthereumConsume.MsgBytes {
				err = utils.ProducerDeliverMessage(msgByte, utils.ToEth, producer)
				if err != nil {
					var msg db.WrapTokenMsg

					if jsonErr := json.Unmarshal(msgByte, &msg); jsonErr != nil {
						logging.Error("Failed to Unmarshal Retry ToEth queue msg", "error:", jsonErr)
					}

					logging.Error("Failed to add msg to kafka ToEth queue (need to do manually), message:", msg, "error:", err)
				}
			}
		}
	}

}
