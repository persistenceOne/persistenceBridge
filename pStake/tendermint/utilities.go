package tendermint

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/relayer/relayer"
	pStakeConfig "github.com/persistenceOne/persistenceCore/pStake/config"
	"github.com/persistenceOne/persistenceCore/pStake/constants"
	"github.com/persistenceOne/persistenceCore/pStake/rest/queries"
)

func GenerateUnsignedTx(chain *relayer.Chain, msgs []sdk.Msg, memo string, timeoutHeight uint64) (signing.Tx, error) {
	ctx := chain.CLIContext(0)

	txf, err := tx.PrepareFactory(ctx, chain.TxFactory(0))
	if err != nil {
		return nil, err
	}

	_, adjusted, err := tx.CalculateGas(ctx.QueryWithData, txf, msgs...)
	if err != nil {
		return nil, err
	}

	txf = txf.WithGas(adjusted).WithMemo(memo).WithTimeoutHeight(timeoutHeight)

	txb, err := tx.BuildUnsignedTx(txf, msgs...)
	if err != nil {
		return nil, err
	}

	return txb.GetTx(), nil
}

// CheckAndGenerateRedelegateMsgs Three possible cases to handle
// 1. Adding n validators and removing m validators, where m = n
// 2. Adding n validators and removing m validators, where m > n
// 3. Adding n validators and removing m validators, where m < n
func CheckAndGenerateRedelegateMsgs() ([]sdk.Msg, error) {
	var stakingMsgs []sdk.Msg
	config := pStakeConfig.GetAppConfiguration()

	mpcValidators, err := queries.GetMPCValidatos(constants.MPCValidatorsURL)
	if err != nil {
		return stakingMsgs, err
	}

	delegations, err := queries.GetDelegations("", config.PStakeAddress.String())
	if err != nil {
		return stakingMsgs, err
	}
	delegationValidators := make([]string, len(delegations.DelegationResponses))
	oldDelegationsMap := map[string]sdk.Int{}
	totalDelegation := sdk.NewInt(0)

	for i, delegationValidator := range delegations.DelegationResponses {
		delegationValidators[i] = delegationValidator.Delegation.ValidatorAddress
		balance, ok := sdk.NewIntFromString(delegationValidator.Balance.Amount)
		if !ok {
			return stakingMsgs, fmt.Errorf("parsing amount from string failed %s", delegationValidator.Balance.Amount)
		}
		totalDelegation = totalDelegation.Add(balance)
		oldDelegationsMap[delegationValidator.Delegation.ValidatorAddress] = balance
	}
	newValidators := difference(mpcValidators.Validators, delegationValidators)
	commonValidators := intersection(mpcValidators.Validators, delegationValidators)
	removedValidators := difference(delegationValidators, mpcValidators.Validators)

	allocateDelegationPerValidator := totalDelegation.Quo(sdk.NewInt(int64(len(mpcValidators.Validators))))

	// What to do with this?
	// leftOver := totalDelegation.Sub(allocateDelegationPerValidator.Mul(sdk.NewInt(int64(len(mpcValidators.Validators)))))

	// This means number of validators removed is greater than added
	if allocateDelegationPerValidator.GT(totalDelegation.Quo(sdk.NewInt(int64(len(delegations.DelegationResponses))))) {
		for k, validator := range append(commonValidators, newValidators...) {
			var transfer sdk.Int
			if k < len(commonValidators) {
				transfer = allocateDelegationPerValidator.Sub(totalDelegation.Quo(sdk.NewInt(int64(len(delegations.DelegationResponses)))))
			} else {
				transfer = allocateDelegationPerValidator
			}
			var srcAmounts []sdk.Int
			var srcValidators []sdk.ValAddress
			for _, removedValidator := range removedValidators {
				removedValidatorDelegationAmt := oldDelegationsMap[removedValidator.String()]
				if removedValidatorDelegationAmt.GT(sdk.ZeroInt()) {
					srcValidators = append(srcValidators, removedValidator)
					if transfer.GT(removedValidatorDelegationAmt) {
						srcAmounts = append(srcAmounts, transfer.Sub(removedValidatorDelegationAmt))
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
				msg := stakingTypes.NewMsgBeginRedelegate(config.PStakeAddress, srcValidator, validator, sdk.NewCoin(config.PStakeDenom, srcAmounts[i]))
				stakingMsgs = append(stakingMsgs, msg)
			}
		}
	}

	// This means number of validators added is greater than removed
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
						srcAmounts = append(srcAmounts, transfer.Sub(removedValidatorDelegationAmt))
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
			if transfer != sdk.ZeroInt() {
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
				msg := stakingTypes.NewMsgBeginRedelegate(config.PStakeAddress, srcValidator, validator, sdk.NewCoin(config.PStakeDenom, srcAmounts[i]))
				stakingMsgs = append(stakingMsgs, msg)
			}
		}
	}

	// If number of mpc validator is same as delegations to validators, it means same numbers of validators has been added and removed.
	if allocateDelegationPerValidator.Equal(totalDelegation.Quo(sdk.NewInt(int64(len(delegations.DelegationResponses))))) {
		for i, newValidator := range newValidators {
			msg := stakingTypes.NewMsgBeginRedelegate(config.PStakeAddress, removedValidators[i], newValidator, sdk.NewCoin(config.PStakeDenom, oldDelegationsMap[removedValidators[i].String()]))
			stakingMsgs = append(stakingMsgs, msg)
		}
	}

	return stakingMsgs, nil
}

// difference returns the elements in `a` that aren't in `b`.
func difference(a, b []string) []sdk.ValAddress {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []sdk.ValAddress
	for _, x := range a {
		if _, found := mb[x]; !found {
			validator, err := sdk.ValAddressFromBech32(x)
			if err != nil {
				panic(err)
			}
			diff = append(diff, validator)
		}
	}
	return diff
}

func intersection(a, b []string) (c []sdk.ValAddress) {
	m := make(map[string]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; ok {
			validator, err := sdk.ValAddressFromBech32(item)
			if err != nil {
				panic(err)
			}
			c = append(c, validator)
		}
	}
	return
}
