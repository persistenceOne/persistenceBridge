/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package handler

import (
	"github.com/Shopify/sarama"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func (m MsgHandler) HandleMsgUnbond(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	config := utils.SaramaConfig()
	producer := utils.NewProducer(configuration.GetAppConfig().Kafka.Brokers, config)

	defer func() {
		err := producer.Close()
		if err != nil {
			logging.Error("failed to close producer in topic: MsgUnbond, error:", err)
		}
	}()

	claimMsgChan := claim.Messages()

	var (
		kafkaMsg *sarama.ConsumerMessage
		ok       bool
	)

ConsumerLoop:
	for {
		select {
		case kafkaMsg, ok = <-claimMsgChan:
			if !ok {
				break ConsumerLoop
			}

			if kafkaMsg == nil {
				return ErrKafkaNilMessage
			}

			err := utils.ProducerDeliverMessage(kafkaMsg.Value, utils.ToTendermint, producer)
			if err != nil {
				// TODO @Puneet return err?? ~ can return, since already logging no logic changes.
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
