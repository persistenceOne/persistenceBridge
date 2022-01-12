/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"github.com/Shopify/sarama"
)

// NewProducer is a producer to send messages to kafka
func NewProducer(kafkaPorts []string, config *sarama.Config) sarama.SyncProducer {
	producer, err := sarama.NewSyncProducer(kafkaPorts, config)
	if err != nil {
		panic(err)
	}

	return producer
}

// ProducerDeliverMessage : delivers messages to kafka
func ProducerDeliverMessage(msgBytes []byte, topic string, producer sarama.SyncProducer) error {
	sendMsg := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msgBytes),
	}

	_, _, err := producer.SendMessage(&sendMsg)

	return err
}

func ProducerDeliverMessages(msgBytes [][]byte, topic string, producer sarama.SyncProducer) error {
	sendMsgs := make([]*sarama.ProducerMessage, len(msgBytes))

	for i, msgByte := range msgBytes {
		sendMsg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(msgByte),
		}

		sendMsgs[i] = sendMsg
	}

	return producer.SendMessages(sendMsgs)
}
