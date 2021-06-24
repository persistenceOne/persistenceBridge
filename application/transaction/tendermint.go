package transaction

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authSigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/relayer/relayer"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/rest/blockchain"
	"github.com/persistenceOne/persistenceBridge/application/rest/casp"
	"log"
)

// Timeout height should be greater than current block height or set it 0 for none.
func GetBytesToSign(chain *relayer.Chain, fromPublicKey cryptotypes.PubKey, msgs []sdk.Msg, memo string, timeoutHeight uint64) ([]byte, client.TxBuilder, tx.Factory, error) {

	from := sdk.AccAddress(fromPublicKey.Address())
	ctx := chain.CLIContext(0).WithFromAddress(from)

	txFactory, err := tx.PrepareFactory(ctx, chain.TxFactory(0))
	if err != nil {
		return []byte{}, nil, txFactory, err
	}

	_, adjusted, err := tx.CalculateGas(ctx.QueryWithData, txFactory, msgs...)
	if err != nil {
		return []byte{}, nil, txFactory, err
	}

	txFactory = txFactory.WithGas(adjusted).WithMemo(memo).WithTimeoutHeight(timeoutHeight)

	txBuilder, err := tx.BuildUnsignedTx(txFactory, msgs...)
	if err != nil {
		return []byte{}, nil, txFactory, err
	}

	signMode := txFactory.SignMode()
	if signMode == signing.SignMode_SIGN_MODE_UNSPECIFIED {
		signMode = ctx.TxConfig.SignModeHandler().DefaultMode()
	}

	signerData := authSigning.SignerData{
		ChainID:       txFactory.ChainID(),
		AccountNumber: txFactory.AccountNumber(),
		Sequence:      txFactory.Sequence(),
	}
	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   fromPublicKey,
		Data:     &sigData,
		Sequence: txFactory.Sequence(),
	}
	if err := txBuilder.SetSignatures(sig); err != nil {
		return []byte{}, txBuilder, txFactory, err
	}

	bytesToSign, err := ctx.TxConfig.SignModeHandler().GetSignBytes(signMode, signerData, txBuilder.GetTx())
	if err != nil {
		return []byte{}, txBuilder, txFactory, err
	}

	return bytesToSign, txBuilder, txFactory, nil
}

// BroadcastMsgs chalk swarm motion broom chapter team guard bracket invest situate circle deny tuition park economy movie subway chase alert popular slogan emerge cricket category
func BroadcastMsgs(chain *relayer.Chain, fromPublicKey cryptotypes.PubKey, sigBytes []byte, txBuilder client.TxBuilder, txFactory tx.Factory) (*sdk.TxResponse, bool, error) {

	from := sdk.AccAddress(fromPublicKey.Address())
	ctx := chain.CLIContext(0).WithFromAddress(from)

	signMode := txFactory.SignMode()
	if signMode == signing.SignMode_SIGN_MODE_UNSPECIFIED {
		signMode = ctx.TxConfig.SignModeHandler().DefaultMode()
	}

	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: sigBytes,
	}
	sig := signing.SignatureV2{
		PubKey:   fromPublicKey,
		Data:     &sigData,
		Sequence: txFactory.Sequence(),
	}

	if err := txBuilder.SetSignatures(sig); err != nil {
		return nil, false, err
	}

	txBytes, err := ctx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, false, err
	}

	res, err := ctx.BroadcastTx(txBytes)
	if err != nil {
		return nil, false, err
	}
	if res.Code != 0 {
		return res, false, nil
	}

	log.Printf("TX HASH: %s, CODE: %d\n", res.TxHash, res.Code)

	return res, true, nil
}

// CheckAndGenerateRedelegateMsgs Three possible cases to handle
// 1. Adding n validators and removing m validators, where m = n
// 2. Adding n validators and removing m validators, where m > n
// 3. Adding n validators and removing m validators, where m < n
func CheckAndGenerateRedelegateMsgs() ([]sdk.Msg, error) {
	var stakingMsgs []sdk.Msg
	config := configuration.GetAppConfiguration()

	mpcValidators, err := casp.GetMPCValidatos(constants.CASP_URL)
	if err != nil {
		return stakingMsgs, err
	}

	delegations, err := blockchain.GetDelegations("", config.PStakeAddress.String())
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
				msg := stakingTypes.NewMsgBeginRedelegate(config.PStakeAddress, srcValidator, validator, sdk.NewCoin(config.PStakeDenom, srcAmounts[i]))
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
