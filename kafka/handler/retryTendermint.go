package handler

import (
	"errors"
	"github.com/Shopify/sarama"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func (m MsgHandler) HandleRetryTendermint(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	config := utils.SaramaConfig()
	producer := utils.NewProducer(configuration.GetAppConfig().Kafka.Brokers, config)
	defer func() {
		err := producer.Close()
		if err != nil {
			logging.Error("failed to close producer in topic RetryTendermint, error:", err)
		}
	}()
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
			if sdk.MsgTypeURL(msg) == constants.MsgSendTypeUrl && !m.WithdrawRewards {
				loop, err := WithdrawRewards(configuration.GetAppConfig().Kafka.ToTendermint.MaxBatchSize-m.Count, m.ProtoCodec, producer, m.Chain)
				if err != nil {
					return err
				}
				m.WithdrawRewards = true
				m.Count = configuration.GetAppConfig().Kafka.ToTendermint.MaxBatchSize - loop
			}

			//TODO remove: This is added as fix for smooth migration.
			if sdk.MsgTypeURL(msg) == constants.MsgDelegateTypeUrl {
				switch txMsg := msg.(type) {
				case *stakingTypes.MsgDelegate:
					if txMsg.Amount.Amount.LTE(sdk.ZeroInt()) {
						session.MarkMessage(kafkaMsg, "")
						continue
					}
				default:
					logging.Fatal("Unexpected type found in topic: EthUnbond")
				}
			}

			err = utils.ProducerDeliverMessage(kafkaMsg.Value, utils.ToTendermint, producer)
			if err != nil {
				//TODO @Puneet return err?? ~ can return, since already logging no logic changes.
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
