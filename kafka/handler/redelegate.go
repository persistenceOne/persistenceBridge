package handler

import (
	"errors"
	"github.com/Shopify/sarama"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/tendermint"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func (m MsgHandler) HandleRelegate(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	config := utils.SaramaConfig()
	producer := utils.NewProducer(configuration.GetAppConfig().Kafka.Brokers, config)
	defer func() {
		err := producer.Close()
		if err != nil {
			logging.Error("failed to close producer in topic: MsgDelegate, error:", err)
		}
	}()

	claimMsgChan := claim.Messages()
	var kafkaMsg *sarama.ConsumerMessage
	var ok bool
	var redelegationSourceAddress sdk.ValAddress
	select {
	case kafkaMsg, ok = <-claimMsgChan:
		if !ok {
			return errors.New("kafka returned error message")
		}
		if kafkaMsg == nil {
			return errors.New("kafka returned nil message")
		}
		redelegationSourceAddress = kafkaMsg.Value
	default:
		return nil
	}

	validatorSet, err := db.GetValidators()
	if err != nil {
		return err
	}
	// query validator delegation
	delegations, err := tendermint.QueryDelegatorDelegations(configuration.GetAppConfig().Tendermint.GetPStakeAddress(), m.Chain)
	if err != nil {
		return err
	}
	totalRedistributeAmount := sdk.ZeroInt()
	for _, delegation := range delegations {
		if delegation.Delegation.ValidatorAddress == redelegationSourceAddress.String() {
			totalRedistributeAmount = delegation.Balance.Amount
		}

	}
	if totalRedistributeAmount.Equal(sdk.ZeroInt()) {
		logging.Info("No Delegations to Redelegate for validator src Address", redelegationSourceAddress.String())
		session.MarkMessage(kafkaMsg, "")
		return nil
	}
	redistributeAmount := totalRedistributeAmount.QuoRaw(int64(len(validatorSet)))
	redistributeChange := totalRedistributeAmount.ModRaw(int64(len(validatorSet)))

	// for loop among validators

	for i, validator := range validatorSet {
		msgRedelegate := &stakingTypes.MsgBeginRedelegate{
			DelegatorAddress:    configuration.GetAppConfig().Tendermint.GetPStakeAddress(),
			ValidatorSrcAddress: redelegationSourceAddress.String(),
			ValidatorDstAddress: validator.Address.String(),
			Amount:              sdk.NewCoin(configuration.GetAppConfig().Tendermint.PStakeDenom, redistributeAmount),
		}
		if i == len(validatorSet)-1 {
			msgRedelegate.Amount.Amount = msgRedelegate.Amount.Amount.Add(redistributeChange)
		}
		msgBytes, err := m.ProtoCodec.MarshalInterface(msgRedelegate)
		if err != nil {
			return err
		}

		err = utils.ProducerDeliverMessage(msgBytes, utils.ToTendermint, producer)
		if err != nil {
			logging.Error("failed to produce message from topic Redelegate to ToTendermint")
			return err
		}
		m.Count++
	}
	session.MarkMessage(kafkaMsg, "")

	return nil
}
