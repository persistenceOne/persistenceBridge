package handler

import (
	"errors"
	"github.com/Shopify/sarama"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/tendermint"
	"log"
)

func (m MsgHandler) HandleEthUnbond(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	saramaConfig := utils.SaramaConfig()
	producer := utils.NewProducer(m.PstakeConfig.Kafka.Brokers, saramaConfig)
	defer func() {
		err := producer.Close()
		if err != nil {
			log.Printf("failed to close producer in topic: %v\n", utils.EthUnbond)
		}
	}()
	var kafkaMsg *sarama.ConsumerMessage

	claimMsgChan := claim.Messages()
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
			case *stakingTypes.MsgUndelegate:
				sum = sum.Add(txMsg.Amount.Amount)
			default:
				log.Printf("Unexpected type found in topic: %v\n", utils.EthUnbond)
			}
		default:
			break ConsumerLoop
		}
	}

	if sum.GT(sdk.NewInt(0)) {
		delegatorDelegations, err := tendermint.QueryDelegatorDelegations(m.Chain.MustGetAddress().String(), m.Chain)
		if err != nil {
			return err
		}
		totalDelegations := TotalDelegations(delegatorDelegations)
		if sum.GT(totalDelegations) {
			return errors.New("Unbondings requested are greater than delegated tokens")
		}
		ratio := sum.ToDec().Quo(totalDelegations.ToDec())
		var unbondings sdk.Int
		var unbondMsgs []*stakingTypes.MsgUndelegate
		for _, delegation := range delegatorDelegations {
			unbondingShare := ratio.Mul(delegation.Balance.Amount.ToDec()).RoundInt()
			if unbondingShare.LT(delegation.Balance.Amount) {
				unbondings = unbondings.Add(unbondingShare)
			} else {
				log.Printf("Incorrect UnbondingShareCalculation: Please Check delegations and unbonding delegations")
			}
			unbondMsg := &stakingTypes.MsgUndelegate{
				DelegatorAddress: m.Chain.MustGetAddress().String(),
				ValidatorAddress: delegation.Delegation.ValidatorAddress,
				Amount: sdk.Coin{
					Denom:  configuration.GetAppConfig().Tendermint.PStakeDenom,
					Amount: sum,
				},
			}
			unbondMsgs = append(unbondMsgs, unbondMsg)
		}

		if unbondings.GT(sum) {
			difference := unbondings.Sub(sum)
		Loop:
			for {
				for _, unbondMsg := range unbondMsgs {
					if unbondMsg.Amount.Amount.GT(sdk.ZeroInt()) {
						unbondMsg.Amount.Amount = unbondMsg.Amount.Amount.Sub(sdk.NewInt(1))
						difference = difference.Sub(sdk.NewInt(1))
					}
					if difference.Equal(sdk.ZeroInt()) {
						break Loop
					}
				}
			}

		}

		for _, unbondMsg := range unbondMsgs {
			msgBytes, err := m.ProtoCodec.MarshalInterface(sdk.Msg(unbondMsg))
			if err != nil {
				return err
			}

			err = utils.ProducerDeliverMessage(msgBytes, utils.MsgUnbond, producer)
			if err != nil {
				log.Printf("failed to produce message from topic %v to %v\n", utils.EthUnbond, utils.MsgUnbond)
				return err
			}
		}
		session.MarkMessage(kafkaMsg, "")
	}
	return nil
}
