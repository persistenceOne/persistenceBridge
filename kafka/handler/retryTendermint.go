/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package handler

import (
	"errors"
	"github.com/Shopify/sarama"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func (m MsgHandler) HandleRetryTendermint(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	config := utils.SaramaConfig()
	producer := utils.NewProducer(configuration.GetAppConfig().Kafka.GetBrokersList(), config)
	defer func() {
		err := producer.Close()
		if err != nil {
			logging.Error("failed to close producer in topic RetryTendermint, error:", err)
		}
	}()
	if m.Count <= 0 {
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
			var msg sdk.Msg
			err := m.ProtoCodec.UnmarshalInterface(kafkaMsg.Value, &msg)
			if err != nil {
				return err
			}
			if msg.Type() == bankTypes.TypeMsgSend && !m.WithdrawRewards {
				validators, err := db.GetValidators()
				if err != nil {
					return err
				}
				loop := configuration.GetAppConfig().Kafka.ToTendermint.MaxBatchSize - m.Count
				if loop <= len(validators) {
					return nil
				}
				loop, err = WithdrawRewards(loop, m.ProtoCodec, producer, m.Chain)
				if err != nil {
					return err
				}
				m.WithdrawRewards = true
				m.Count = configuration.GetAppConfig().Kafka.ToTendermint.MaxBatchSize - loop
				if !checkCount(m.Count, configuration.GetAppConfig().Kafka.ToTendermint.MaxBatchSize) {
					break ConsumerLoop
				}
			}

			err = utils.ProducerDeliverMessage(kafkaMsg.Value, utils.ToTendermint, producer)
			if err != nil {
				logging.Error("failed to produce from: RetryTendermint to: ToTendermint, error:", err)
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
