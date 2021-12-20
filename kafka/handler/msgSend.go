/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package handler

import (
	"github.com/Shopify/sarama"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func (m MsgHandler) HandleMsgSend(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	config := utils.SaramaConfig()
	producer := utils.NewProducer(configuration.GetAppConfig().Kafka.Brokers, config)

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

				if !m.WithdrawRewards {
					loop, err = WithdrawRewards(loop, m.ProtoCodec, producer, m.Chain)
					if err != nil {
						return err
					}

					m.WithdrawRewards = true
				}

				err = utils.ProducerDeliverMessage(kafkaMsg.Value, utils.ToTendermint, producer)
				if err != nil {
					// TODO @Puneet return err?? ~ can return, since already logging no logic changes.
					logging.Error("failed to produce from: MsgSend to: ToTendermint")

					break ConsumerLoop
				}

				session.MarkMessage(kafkaMsg, "")

				loop--
				if loop == 0 {
					break ConsumerLoop
				}
			default:
				break ConsumerLoop
			}
		}
	}

	return nil
}
