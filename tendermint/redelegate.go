package tendermint

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/rest/blockchain"
)

// CheckAndGenerateRedelegateMsgs Three possible cases to handle
// 1. Adding n validators and removing m validators, where m = n
// 2. Adding n validators and removing m validators, where m > n
// 3. Adding n validators and removing m validators, where m < n
func CheckAndGenerateRedelegateMsgs() ([]sdk.Msg, error) {
	var stakingMsgs []sdk.Msg
	config := configuration.GetAppConfig()

	validators, err := db.GetValidators()
	if err != nil {
		return stakingMsgs, err
	}

	delegations, err := blockchain.GetDelegations("", config.Tendermint.PStakeAddress.String())
	if err != nil {
		return stakingMsgs, err
	}
	delegationValidators := make([]sdk.ValAddress, len(delegations.DelegationResponses))
	oldDelegationsMap := map[string]sdk.Int{}
	totalDelegation := sdk.NewInt(0)

	for i, delegationValidator := range delegations.DelegationResponses {
		delegationValidators[i], err = sdk.ValAddressFromBech32(delegationValidator.Delegation.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		balance, ok := sdk.NewIntFromString(delegationValidator.Balance.Amount)
		if !ok {
			return stakingMsgs, fmt.Errorf("parsing amount from string failed %s", delegationValidator.Balance.Amount)
		}
		totalDelegation = totalDelegation.Add(balance)
		oldDelegationsMap[delegationValidator.Delegation.ValidatorAddress] = balance
	}
	newValidators := difference(validators, delegationValidators)
	commonValidators := intersection(validators, delegationValidators)
	removedValidators := difference(delegationValidators, validators)

	allocateDelegationPerValidator := totalDelegation.Quo(sdk.NewInt(int64(len(validators))))

	// What to do with this?
	// leftOver := totalDelegation.Sub(allocateDelegationPerValidator.Mul(sdk.NewInt(int64(len(validators.Validators)))))

	// This means number of validators removed is greater than added. Here we have transfer to common and new validators subtracting from removed validators.
	if allocateDelegationPerValidator.GT(totalDelegation.Quo(sdk.NewInt(int64(len(delegations.DelegationResponses))))) {
		for k, validator := range append(commonValidators, newValidators...) {
			var transfer sdk.Int
			if k < len(commonValidators) {
				// transferring to common validator
				transfer = allocateDelegationPerValidator.Sub(totalDelegation.Quo(sdk.NewInt(int64(len(delegations.DelegationResponses)))))
			} else {
				// transferring to new validator
				transfer = allocateDelegationPerValidator
			}
			var srcAmounts []sdk.Int
			var srcValidators []sdk.ValAddress
			for _, removedValidator := range removedValidators {
				removedValidatorDelegationAmt := oldDelegationsMap[removedValidator.String()]
				if removedValidatorDelegationAmt.GT(sdk.ZeroInt()) {
					srcValidators = append(srcValidators, removedValidator)
					if transfer.GT(removedValidatorDelegationAmt) {
						srcAmounts = append(srcAmounts, removedValidatorDelegationAmt)
						oldDelegationsMap[removedValidator.String()] = sdk.ZeroInt()
						transfer = transfer.Sub(removedValidatorDelegationAmt)
					} else {
						srcAmounts = append(srcAmounts, transfer)
						oldDelegationsMap[removedValidator.String()] = removedValidatorDelegationAmt.Sub(transfer)
						transfer = sdk.ZeroInt()
						break
					}
				}
			}
			if len(srcAmounts) != len(srcValidators) {
				panic("invalid code")
			}
			for i, srcValidator := range srcValidators {
				msg := stakingTypes.NewMsgBeginRedelegate(config.Tendermint.PStakeAddress, srcValidator, validator, sdk.NewCoin(config.Tendermint.PStakeDenom, srcAmounts[i]))
				stakingMsgs = append(stakingMsgs, msg)
			}
		}
	}

	// This means number of validators added is greater than removed. Here we have transfer only to new validators subtracting from removed and common validators.
	if allocateDelegationPerValidator.LT(totalDelegation.Quo(sdk.NewInt(int64(len(delegations.DelegationResponses))))) {
		for _, validator := range newValidators {
			transfer := allocateDelegationPerValidator
			var srcAmounts []sdk.Int
			var srcValidators []sdk.ValAddress
			for _, removedValidator := range removedValidators {
				removedValidatorDelegationAmt := oldDelegationsMap[removedValidator.String()]
				if removedValidatorDelegationAmt.GT(sdk.ZeroInt()) {
					srcValidators = append(srcValidators, removedValidator)
					if transfer.GT(removedValidatorDelegationAmt) {
						srcAmounts = append(srcAmounts, removedValidatorDelegationAmt)
						oldDelegationsMap[removedValidator.String()] = sdk.ZeroInt()
						transfer = transfer.Sub(removedValidatorDelegationAmt)
					} else {
						srcAmounts = append(srcAmounts, transfer)
						oldDelegationsMap[removedValidator.String()] = removedValidatorDelegationAmt.Sub(transfer)
						transfer = sdk.ZeroInt()
						break
					}
				}
			}
			if !transfer.Equal(sdk.ZeroInt()) {
				amount := transfer.Quo(sdk.NewInt(int64(len(commonValidators)))) // can lead to non zero leftover
				for _, commonValidator := range commonValidators {
					srcValidators = append(srcValidators, commonValidator)
					srcAmounts = append(srcAmounts, amount)
				}
			}
			if len(srcAmounts) != len(srcValidators) {
				panic("invalid code")
			}
			for i, srcValidator := range srcValidators {
				msg := stakingTypes.NewMsgBeginRedelegate(config.Tendermint.PStakeAddress, srcValidator, validator, sdk.NewCoin(config.Tendermint.PStakeDenom, srcAmounts[i]))
				stakingMsgs = append(stakingMsgs, msg)
			}
		}
	}

	// If number of mpc validator is same as delegations to validators, it means same numbers of validators has been added and removed.
	if allocateDelegationPerValidator.Equal(totalDelegation.Quo(sdk.NewInt(int64(len(delegations.DelegationResponses))))) {
		for i, newValidator := range newValidators {
			msg := stakingTypes.NewMsgBeginRedelegate(config.Tendermint.PStakeAddress, removedValidators[i], newValidator, sdk.NewCoin(config.Tendermint.PStakeDenom, oldDelegationsMap[removedValidators[i].String()]))
			stakingMsgs = append(stakingMsgs, msg)
		}
	}

	return stakingMsgs, nil
}

// difference returns the elements in `a` that aren't in `b`.
func difference(a, b []sdk.ValAddress) []sdk.ValAddress {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x.String()] = struct{}{}
	}
	var diff []sdk.ValAddress
	for _, x := range a {
		if _, found := mb[x.String()]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func intersection(a, b []sdk.ValAddress) (c []sdk.ValAddress) {
	m := make(map[string]bool)

	for _, item := range a {
		m[item.String()] = true
	}

	for _, item := range b {
		if _, ok := m[item.String()]; ok {
			c = append(c, item)
		}
	}
	return
}
