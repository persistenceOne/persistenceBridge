package handler

import (
	"errors"
	"github.com/Shopify/sarama"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"log"
)

func (m MsgHandler) HandleMsgDelegate(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	config := utils.SaramaConfig()
	producer := utils.NewProducer(m.PstakeConfig.Kafka.Brokers, config)
	defer func() {
		err := producer.Close()
		if err != nil {
			log.Printf("failed to close producer in topic: %v\n", utils.MsgDelegate)
		}
	}()

	claimMsgChan := claim.Messages()
	var kafkaMsg *sarama.ConsumerMessage
	var ok bool
	var sum = sdk.NewInt(0)
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
				log.Printf("proto failed to unmarshal\n")
			}
			switch txMsg := msg.(type) {
			case *stakingTypes.MsgDelegate:
				sum = sum.Add(txMsg.Amount.Amount)
			default:
				log.Printf("Unexpected type found in topic: %v\n", utils.EthUnbond)
			}
		default:
			break ConsumerLoop
		}
	}

	if sum.GT(sdk.NewInt(0)) {
		// TODO consider multiple validators
		delegateMsg := &stakingTypes.MsgDelegate{
			DelegatorAddress: m.Chain.MustGetAddress().String(),
			ValidatorAddress: constants2.Validator1.String(),
			Amount: sdk.Coin{
				Denom:  configuration.GetAppConfiguration().Tendermint.PStakeDenom,
				Amount: sum,
			},
		}
		msgBytes, err := m.ProtoCodec.MarshalInterface(sdk.Msg(delegateMsg))
		if err != nil {
			return err
		}

		err = utils.ProducerDeliverMessage(msgBytes, utils.ToTendermint, producer)
		if err != nil {
			log.Printf("failed to produce message from topic %v to %v\n", utils.MsgDelegate, utils.ToTendermint)
			return err
		}
		session.MarkMessage(kafkaMsg, "")
		m.Count++
	}
	return nil
}
