/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package handler

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func (m MsgHandler) HandleMsgSend(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	config := utils.SaramaConfig()
	producer := utils.NewProducer(configuration.GetAppConfig().Kafka.GetBrokersList(), config)
	defer func() {
		err := producer.Close()
		if err != nil {
			logging.Error("failed to close producer in topic: MsgSend, error:", err)
		}
	}()
	validators, err := db.GetValidators()
	if err != nil {
		return err
	}

	loop := configuration.GetAppConfig().Kafka.ToTendermint.MaxBatchSize - m.Count
	if loop <= len(validators) || m.WithdrawRewards {
		return nil
	}

	if loop > 0 {
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

				if !m.WithdrawRewards {
					loop, err = WithdrawRewards(loop, m.ProtoCodec, producer, m.Chain)
					if err != nil {
						return err
					}
					m.Count = m.Count + len(validators)
					if !checkCount(m.Count, configuration.GetAppConfig().Kafka.ToTendermint.MaxBatchSize) {
						break ConsumerLoop
					}
					m.WithdrawRewards = true
				}
				err := utils.ProducerDeliverMessage(kafkaMsg.Value, utils.ToTendermint, producer)
				if err != nil {
					logging.Error("failed to produce from: MsgSend to: ToTendermint")
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
	}
	return nil
}
