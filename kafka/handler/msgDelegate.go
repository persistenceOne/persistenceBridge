/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package handler

import (
	"github.com/Shopify/sarama"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func (m *MsgHandler) HandleMsgDelegate(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	config := utils.SaramaConfig()
	producer := utils.NewProducer(configuration.GetAppConfig().Kafka.Brokers, config)

	defer func() {
		err := producer.Close()
		if err != nil {
			logging.Error("failed to close producer in topic: MsgDelegate, bridgeErr:", err)
		}
	}()

	claimMsgChan := claim.Messages()

	var (
		kafkaMsg *sarama.ConsumerMessage
		ok       bool
	)

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
			case *stakingTypes.MsgDelegate:
				sum = sum.Add(txMsg.Amount.Amount)
			default:
				logging.Error("Unexpected type found in topic: EthUnbond")
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
			logging.Error("Delegate transaction number is higher than slots available, probably increase to tendermint MaxBatchSize")

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

			if !delegateMsg.Amount.Amount.LTE(sdk.ZeroInt()) {
				var msgBytes []byte

				msgBytes, err = m.ProtoCodec.MarshalInterface(delegateMsg)
				if err != nil {
					return err
				}

				err = utils.ProducerDeliverMessage(msgBytes, utils.ToTendermint, producer)
				if err != nil {
					logging.Error("failed to produce message from: MsgDelegate to ToTendermint")

					return err
				}

				m.Count++
			}
		}

		session.MarkMessage(kafkaMsg, "")
	}

	return nil
}
