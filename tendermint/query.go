package tendermint

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/relayer/relayer"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
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
		logging.Error("Delegator's validator not found, Error:", err)
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
		logging.Error("Delegator validator delegation not found, Error:", err)
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
		logging.Info("Delegator delegations not found, Error:", err)
		return nil, err
	}
	return stakingRes.DelegationResponses, err
}

func QueryValidatorStatus(validatorAddress sdk.ValAddress, chain *relayer.Chain) (stakingTypes.QueryValidatorResponse, error) {
	stakingClient := stakingTypes.NewQueryClient(chain.CLIContext(0))
	stakingResponse, err := stakingClient.Validator(context.Background(), &stakingTypes.QueryValidatorRequest{ValidatorAddr: validatorAddress.String()})
	if err != nil {
		logging.Info("Validator's staking response not found, Error:", err)
		return stakingTypes.QueryValidatorResponse{}, err
	}
	return *stakingResponse, err
}

func QueryAllValidators(chain *relayer.Chain) ([]rpc.ValidatorOutput, error) {
	height, _ := chain.QueryLatestHeight()
	var (
		page    = 1
		perPage = 10000
	)
	valOutput, err := rpc.GetValidators(chain.CLIContext(0), &height, &page, &perPage)
	if err != nil {
		logging.Info("List of Validators not found, Error:", err)
		return nil, err
	}
	return valOutput.Validators, err
}

func QuerySlashingSigningInfo(consAddress sdk.ConsAddress, chain *relayer.Chain) (slashingTypes.QuerySigningInfoResponse, error) {
	slashingClient := slashingTypes.NewQueryClient(chain.CLIContext(0))
	slashingResponse, err := slashingClient.SigningInfo(context.Background(), &slashingTypes.QuerySigningInfoRequest{ConsAddress: consAddress.String()})
	if err != nil {
		logging.Info("Validator's signing info not found, Error:", err)
		return slashingTypes.QuerySigningInfoResponse{}, err
	}
	return *slashingResponse, err
}
