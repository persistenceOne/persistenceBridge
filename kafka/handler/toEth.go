/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package handler

import (
	"encoding/json"
	"time"

	"github.com/Shopify/sarama"

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

	err := SendBatchToEth(kafkaMsgs, m)
	if err != nil {
		return err
	}

	session.MarkMessage(kafkaMsg, "")

	return nil
}

func convertKafkaMsgsToEthMsg(kafkaMsgs []sarama.ConsumerMessage) ([]outgoingtx.WrapTokenMsg, error) {
	msgs := make([]outgoingtx.WrapTokenMsg, len(kafkaMsgs))

	for i := range kafkaMsgs {
		var msg outgoingtx.WrapTokenMsg

		err := json.Unmarshal(kafkaMsgs[i].Value, &msg)
		if err != nil {
			return nil, err
		}

		msgs[i] = msg
	}

	return msgs, nil
}

// SendBatchToEth : Handling of msgSend
func SendBatchToEth(kafkaMsgs []sarama.ConsumerMessage, handler MsgHandler) error {
	msgs, err := convertKafkaMsgsToEthMsg(kafkaMsgs)
	if err != nil {
		return err
	}

	logging.Info("batched messages to send to ETH:", msgs)

	hash, err := outgoingtx.EthereumWrapToken(handler.EthClient, msgs)
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

		for i := range kafkaMsgs {
			err = utils.ProducerDeliverMessage(kafkaMsgs[i].Value, utils.ToEth, producer)
			if err != nil {
				logging.Error("Failed to add msg to kafka queue, message:", msgs[i], "error:", err)
				// TODO @Puneet continue or return? ~ Log (ALERT) and continue, need to manually do the failed ones.
			}
		}

		return err
	}

	err = db.SetOutgoingEthereumTx(db.NewOutgoingETHTransaction(hash, msgs))
	if err != nil {
		logging.Fatal(err)
	}

	logging.Info("Broadcast Eth Tx hash:", hash.String())

	return nil
}
