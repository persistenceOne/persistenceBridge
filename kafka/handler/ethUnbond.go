/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package handler

import (
	"fmt"

	"github.com/Shopify/sarama"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/tendermint"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func (m MsgHandler) HandleEthUnbond(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	saramaConfig := utils.SaramaConfig()
	producer := utils.NewProducer(configuration.GetAppConfig().Kafka.Brokers, saramaConfig)

	defer func() {
		err := producer.Close()
		if err != nil {
			logging.Error("failed to close producer in topic: EthUnbond, bridgeErr:", err)
		}
	}()

	var (
		kafkaMsg *sarama.ConsumerMessage
		ok       bool
	)

	claimMsgChan := claim.Messages()
	sum := sdk.ZeroInt()

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

			var msg sdk.Msg
			err := m.ProtoCodec.UnmarshalInterface(kafkaMsg.Value, &msg)
			if err != nil {
				return err
			}

			switch txMsg := msg.(type) {
			case *stakingTypes.MsgUndelegate:
				sum = sum.Add(txMsg.Amount.Amount)
			default:
				logging.Error("Unexpected type found in topic: EthUnbond")
			}
		default:
			break ConsumerLoop
		}
	}

	if sum.GT(sdk.ZeroInt()) {
		delegatorDelegations, err := tendermint.QueryDelegatorDelegations(configuration.GetAppConfig().Tendermint.GetPStakeAddress(), m.Chain)
		if err != nil {
			return err
		}

		totalDelegations := TotalDelegations(delegatorDelegations)
		if sum.GT(totalDelegations) {
			return fmt.Errorf("%w: unbondings requested %s , delegated %s", ErrTooFewDelegatedTokens, sum, totalDelegations)
		}

		ratio := sum.ToDec().Quo(totalDelegations.ToDec())
		unbondings := sdk.ZeroInt()

		var unbondMsgs []*stakingTypes.MsgUndelegate

		for _, delegation := range delegatorDelegations {
			unbondingShare := ratio.Mul(delegation.Balance.Amount.ToDec()).RoundInt()

			if unbondingShare.LT(delegation.Balance.Amount) {
				unbondings = unbondings.Add(unbondingShare)
			} else {
				logging.Error("Incorrect UnbondingShareCalculation: Please Check delegations and unbonding delegations")
			}

			unbondMsg := &stakingTypes.MsgUndelegate{
				DelegatorAddress: configuration.GetAppConfig().Tendermint.GetPStakeAddress(),
				ValidatorAddress: delegation.Delegation.ValidatorAddress,
				Amount: sdk.Coin{
					Denom:  configuration.GetAppConfig().Tendermint.PStakeDenom,
					Amount: unbondingShare,
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
			if unbondMsg.Amount.Amount.LTE(sdk.ZeroInt()) {
				continue
			}

			var msgBytes []byte

			msgBytes, err = m.ProtoCodec.MarshalInterface(unbondMsg)
			if err != nil {
				return err
			}

			err = utils.ProducerDeliverMessage(msgBytes, utils.MsgUnbond, producer)
			if err != nil {
				logging.Error("failed to produce message from: EthUnbond to: MsgUnbond")

				return err
			}
		}

		session.MarkMessage(kafkaMsg, "")
	}

	return nil
}
