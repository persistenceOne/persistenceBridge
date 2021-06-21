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

	messagesLength := len(claim.Messages())
	loop := messagesLength
	if m.PstakeConfig.Kafka.ToTendermint.BatchSize-m.Count <= 0 {
		return nil
	}
	if messagesLength > m.PstakeConfig.Kafka.ToTendermint.BatchSize-m.Count {
		loop = m.PstakeConfig.Kafka.ToTendermint.BatchSize - m.Count
	}
	if messagesLength > 0 {
		var msgs [][]byte
		// TODO add msg withdraw rewards from multiple validators.
		if tendermint.AddressIsDelegatorToValidator(m.Chain.MustGetAddress().String(), constants2.Validator1.String(), m.Chain) {
			withdrawRewardsMsg := &distributionTypes.MsgWithdrawDelegatorReward{
				DelegatorAddress: m.Chain.MustGetAddress().String(),
				ValidatorAddress: constants2.Validator1.String(),
			}
			withdrawRewardsMsgBytes, err := m.ProtoCodec.MarshalInterface(sdk.Msg(withdrawRewardsMsg))
			if err != nil {
				log.Printf("Failed to Marshal WithdrawMessage: Error: %v\n", err)
			} else {
				msgs = append(msgs, withdrawRewardsMsgBytes)
				loop = loop - 1
			}
		}

		var kafkaMsg *sarama.ConsumerMessage
		for i := 0; i < loop; i++ {
			kafkaMsg = <-claim.Messages()
			if kafkaMsg == nil {
				return errors.New("kafka returned nil message")
			}
			msgs = append(msgs, kafkaMsg.Value)
		}
		if len(msgs) > 0 {
			err := utils.ProducerDeliverMessages(msgs, utils.ToTendermint, producer)
			session.MarkMessage(kafkaMsg, "")
			if err != nil {
				log.Printf("error in handler for topic %v, failed to produce to queue\n", utils.MsgSend)
				return err
			}
		}
	}
	m.Count += loop
	return nil
}
