package handler

import (
	"errors"
	"github.com/Shopify/sarama"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"log"
)

func (m MsgHandler) HandleMsgDelegate(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	config := utils.SaramaConfig()
	producer := utils.NewProducer(configuration.GetAppConfig().Kafka.Brokers, config)
	defer func() {
		err := producer.Close()
		if err != nil {
			log.Printf("failed to close producer in topic: %v\n", utils.MsgDelegate)
		}
	}()

	claimMsgChan := claim.Messages()
	var kafkaMsg *sarama.ConsumerMessage
	var ok bool
	sum := sdk.ZeroInt()
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
		validators, err := db.GetValidators()
		if err != nil {
			return err
		}
		if configuration.GetAppConfig().Kafka.ToTendermint.MaxBatchSize-m.Count < len(validators) {
			log.Printf("Delegate transaction number is higher than slots available, probably increase to tendermint MaxBatchSize")
			return nil
		}
		delegationAmount := sum.QuoRaw(int64(len(validators)))
		delegationChange := sum.ModRaw(int64(len(validators)))

		for i, validator := range validators {
			delegateMsg := &stakingTypes.MsgDelegate{
				DelegatorAddress: configuration.GetAppConfig().Tendermint.GetPStakeAddress(),
				ValidatorAddress: validator.Address.String(),
				Amount: sdk.Coin{
					Denom:  configuration.GetAppConfig().Tendermint.PStakeDenom,
					Amount: delegationAmount,
				},
			}
			if i == len(validators)-1 {
				delegateMsg.Amount.Amount = delegateMsg.Amount.Amount.Add(delegationChange)
			}
			msgBytes, err := m.ProtoCodec.MarshalInterface(delegateMsg)
			if err != nil {
				return err
			}

			err = utils.ProducerDeliverMessage(msgBytes, utils.ToTendermint, producer)
			if err != nil {
				log.Printf("failed to produce message from topic %v to %v\n", utils.MsgDelegate, utils.ToTendermint)
				return err
			}
			m.Count++
		}
		session.MarkMessage(kafkaMsg, "")
	}
	return nil
}
