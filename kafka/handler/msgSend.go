package handler

import (
	"errors"
	"github.com/Shopify/sarama"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/tendermint"
	"log"
)

func (m MsgHandler) HandleMsgSend(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	config := utils.SaramaConfig()
	producer := utils.NewProducer(m.PstakeConfig.Kafka.Brokers, config)
	defer func() {
		err := producer.Close()
		if err != nil {
			log.Printf("failed to close producer in topic: %v\n", utils.MsgSend)
		}
	}()
	loop := m.PstakeConfig.Kafka.ToTendermint.MaxBatchSize - m.Count
	if loop <= 0 {
		return nil
	}

	// TODO add msg withdraw rewards from multiple validators.
	if tendermint.AddressIsDelegatorToValidator(m.Chain.MustGetAddress().String(), constants2.Validator1.String(), m.Chain) {
		withdrawRewardsMsg := &distributionTypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: m.Chain.MustGetAddress().String(),
			ValidatorAddress: constants2.Validator1.String(),
		}
		withdrawRewardsMsgBytes, err := m.ProtoCodec.MarshalInterface(sdk.Msg(withdrawRewardsMsg))
		if err != nil {
			log.Printf("Failed to Marshal WithdrawMessage: Error: %v\n", err)
			return err
		} else {
			err2 := utils.ProducerDeliverMessage(withdrawRewardsMsgBytes, utils.ToTendermint, producer)
			if err2 != nil {
				log.Printf("error in handler for topic %v, failed to produce to queue\n", utils.MsgSend)
				return err2
			}
			loop = loop - 1
		}
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
				err := utils.ProducerDeliverMessage(kafkaMsg.Value, utils.ToTendermint, producer)
				if err != nil {
					log.Printf("failed to produce from %v to :%v", utils.MsgSend, utils.ToTendermint)
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