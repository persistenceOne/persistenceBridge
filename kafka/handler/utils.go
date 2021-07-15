package handler

import (
	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/relayer/relayer"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/tendermint"
	"log"
)

func contains(s []sdk.ValAddress, e sdk.ValAddress) bool {
	for _, a := range s {
		if a.String() == e.String() {
			return true
		}
	}
	return false
}

func ValidatorsInDelegations(delegationResponses stakingTypes.DelegationResponses) []sdk.ValAddress {
	var validators []sdk.ValAddress
	for _, delegation := range delegationResponses {
		validators = append(validators, delegation.Delegation.GetValidatorAddr())
	}
	return validators
}

func TotalDelegations(delegationResponses stakingTypes.DelegationResponses) sdk.Int {
	sum := sdk.ZeroInt()
	for _, delegation := range delegationResponses {
		sum = sum.Add(delegation.Balance.Amount)
	}
	return sum
}

func checkCount(currentCount, maxCount int) bool {
	if currentCount < maxCount {
		return true
	}
	return false
}

func WithdrawRewards(loop int, protoCodec *codec.ProtoCodec, producer sarama.SyncProducer, chain *relayer.Chain) (int, error) {
	validators, err := db.GetValidators()
	if err != nil {
		return loop, err
	}
	delegatorDelegations, err := tendermint.QueryDelegatorDelegations(configuration.GetAppConfig().Tendermint.GetPStakeAddress(), chain)
	if err != nil {
		return loop, err
	}
	delegatorValidators := ValidatorsInDelegations(delegatorDelegations)
	for _, validator := range validators {
		if contains(delegatorValidators, validator.Address) {
			withdrawRewardsMsg := &distributionTypes.MsgWithdrawDelegatorReward{
				DelegatorAddress: configuration.GetAppConfig().Tendermint.GetPStakeAddress(),
				ValidatorAddress: validator.Address.String(),
			}
			withdrawRewardsMsgBytes, err := protoCodec.MarshalInterface(withdrawRewardsMsg)
			if err != nil {
				log.Printf("Failed to Marshal WithdrawMessage: Error: %v\n", err)
				return loop, err
			} else {
				err2 := utils.ProducerDeliverMessage(withdrawRewardsMsgBytes, utils.ToTendermint, producer)
				if err2 != nil {
					log.Printf("error in handler for topic %v, failed to produce to queue\n", utils.MsgSend)
					return loop, err2
				}
				loop = loop - 1
				if loop == 0 {
					return loop, nil
				}
			}
		}
	}
	return loop, nil
}
