package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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
	var sum sdk.Int
	for _, delegation := range delegationResponses {
		sum = sum.Add(delegation.Balance.Amount)
	}
	return sum
}
