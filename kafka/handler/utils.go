/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package handler

import (
	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/relayer/relayer"
	"github.com/dgraph-io/badger/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/tendermint"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
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
	validators := make([]sdk.ValAddress, len(delegationResponses))

	for i, delegation := range delegationResponses {
		validators[i] = delegation.Delegation.GetValidatorAddr()
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
	return currentCount < maxCount
}

func WithdrawRewards(loop int, protoCodec *codec.ProtoCodec, producer sarama.SyncProducer, chain *relayer.Chain, database *badger.DB) (int, error) {
	validators, err := db.GetValidators(database)
	if err != nil {
		return loop, err
	}

	delegatorDelegations, err := tendermint.QueryDelegatorDelegations(configuration.GetAppConfig().Tendermint.GetPStakeAddress(), chain)
	if err != nil {
		errStatus, ok := status.FromError(err)
		if ok && errStatus.Code() == codes.NotFound {
			return loop, nil
		}

		return loop, err
	}

	delegatorValidators := ValidatorsInDelegations(delegatorDelegations)

	for _, validator := range validators {
		if !contains(delegatorValidators, validator.Address) {
			continue
		}

		withdrawRewardsMsg := &distributionTypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: configuration.GetAppConfig().Tendermint.GetPStakeAddress(),
			ValidatorAddress: validator.Address.String(),
		}

		withdrawRewardsMsgBytes, err := protoCodec.MarshalInterface(withdrawRewardsMsg)
		if err != nil {
			logging.Error("Failed to Marshal WithdrawMessage, Error:", err)

			return loop, err
		}

		err = utils.ProducerDeliverMessage(withdrawRewardsMsgBytes, utils.ToTendermint, producer)
		if err != nil {
			logging.Error("failed to produce withdrawRewards to queue ToTendermint")

			return loop, err
		}

		loop--
		if loop == 0 {
			return loop, nil
		}
	}

	return loop, nil
}
