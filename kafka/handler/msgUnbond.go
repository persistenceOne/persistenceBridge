/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package handler

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func (m MsgHandler) HandleMsgUnbond(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	config := utils.SaramaConfig()
	producer := utils.NewProducer(configuration.GetAppConfig().Kafka.GetBrokersList(), config)
	defer func() {
		err := producer.Close()
		if err != nil {
			logging.Error("failed to close producer in topic: MsgUnbond, error:", err)
		}
	}()

	if !checkCount(m.Count, configuration.GetAppConfig().Kafka.ToTendermint.MaxBatchSize) {
		return nil
	}

	claimMsgChan := claim.Messages()
	var kafkaMsg *sarama.ConsumerMessage
	var ok bool
ConsumerLoop:
	for {
		select {
		case kafkaMsg, ok = <-claimMsgChan:
			if !ok {
				break ConsumerLoop
			}
			if kafkaMsg == nil {
				return errors.New("kafka returned nil message")
			}
			err := utils.ProducerDeliverMessage(kafkaMsg.Value, utils.ToTendermint, producer)
			if err != nil {
				logging.Error("failed to produce from MsgUnbond to ToTendermint, error:", err)
				break ConsumerLoop
			}
			session.MarkMessage(kafkaMsg, "")
			m.Count++
			if !checkCount(m.Count, configuration.GetAppConfig().Kafka.ToTendermint.MaxBatchSize) {
				break ConsumerLoop
			}
		default:
			break ConsumerLoop
		}
	}

	return nil
}
