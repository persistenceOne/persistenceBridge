package tendermint

import (
	"context"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/relayer/relayer"
	"log"
)

func AddressIsDelegatorToValidator(delegatorAddress, validatorAddress string, chain *relayer.Chain) bool {
	_, err := QueryValidatorDelegator(delegatorAddress, validatorAddress, chain)
	if err != nil {
		return false
	}
	return true
}

func QueryValidatorDelegator(delegatorAddress, validatorAddress string, chain *relayer.Chain) (stakingTypes.Validator, error) {
	stakingClient := stakingTypes.NewQueryClient(chain.CLIContext(0))
	stakingRes, err := stakingClient.DelegatorValidator(context.Background(), &stakingTypes.QueryDelegatorValidatorRequest{
		DelegatorAddr: delegatorAddress,
		ValidatorAddr: validatorAddress,
	})
	if err != nil {
		log.Printf("Delegator delegations not found, Error: %v\n", err)
		return stakingTypes.Validator{}, err
	}
	return stakingRes.GetValidator(), err
}

func QueryDelegatorValidatorDelegations(delegatorAddress, validatorAddress string, chain *relayer.Chain) (stakingTypes.DelegationResponse, error) {
	stakingClient := stakingTypes.NewQueryClient(chain.CLIContext(0))
	stakingRes, err := stakingClient.Delegation(context.Background(), &stakingTypes.QueryDelegationRequest{
		DelegatorAddr: delegatorAddress,
		ValidatorAddr: validatorAddress,
	})
	if err != nil {
		log.Printf("Delegator delegations not found, Error: %v\n", err)
		return stakingTypes.DelegationResponse{}, err
	}
	return *stakingRes.DelegationResponse, err
}

func QueryDelegatorDelegations(delegatorAddress string, chain *relayer.Chain) (stakingTypes.DelegationResponses, error) {
	stakingClient := stakingTypes.NewQueryClient(chain.CLIContext(0))
	stakingRes, err := stakingClient.DelegatorDelegations(context.Background(), &stakingTypes.QueryDelegatorDelegationsRequest{
		DelegatorAddr: delegatorAddress,
	})
	if err != nil {
		log.Printf("Delegator delegations not found, Error: %v\n", err)
		return nil, err
	}
	return stakingRes.DelegationResponses, err
}
