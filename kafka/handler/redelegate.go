package handler

import (
	"errors"
	"github.com/Shopify/sarama"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/tendermint"
	"log"
)

func (m MsgHandler) HandleRelegate(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
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
	delegations, err := tendermint.QueryDelegatorDelegations(m.Chain.MustGetAddress().String(), m.Chain)
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
		log.Printf("No Delegations to Redelegate for validator src Address %v", redelegationSourceAddress.String())
		return nil
	}
	redistributeAmount := totalRedistributeAmount.QuoRaw(int64(len(validatorSet)))
	redistributeChange := totalRedistributeAmount.ModRaw(int64(len(validatorSet)))

	// for loop among validators

	for i, validator := range validatorSet {
		msgRedelegate := stakingTypes.NewMsgBeginRedelegate(m.Chain.MustGetAddress(), redelegationSourceAddress, validator, sdk.NewCoin(m.PstakeConfig.Tendermint.PStakeDenom, redistributeAmount))
		if i == len(validatorSet)-1 {
			msgRedelegate.Amount.Amount = msgRedelegate.Amount.Amount.Add(redistributeChange)
		}
		msgBytes, err := m.ProtoCodec.MarshalInterface(sdk.Msg(msgRedelegate))
		if err != nil {
			return err
		}

		err = utils.ProducerDeliverMessage(msgBytes, utils.ToTendermint, producer)
		if err != nil {
			log.Printf("failed to produce message from topic %v to %v\n", utils.Redelegate, utils.ToTendermint)
			return err
		}
		m.Count++
	}
	session.MarkMessage(kafkaMsg, "")

	return nil
}
