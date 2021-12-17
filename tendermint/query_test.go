/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package tendermint

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceBridge/application/constants"
)

func TestAddressIsDelegatorToValidator(t *testing.T) {
	chain, _ := InitializeAndStartChain("336h", constants.DefaultPBridgeHome)
	delegatorToValidator := AddressIsDelegatorToValidator("cosmos1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d7lzl74n", "cosmosvaloper1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d76ktteq", chain)
	require.Equal(t, true, delegatorToValidator)
}

func TestQueryDelegatorDelegations(t *testing.T) {
	chain, _ := InitializeAndStartChain("336h", constants.DefaultPBridgeHome)
	queryResponse, err := QueryDelegatorDelegations("cosmos1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d7lzl74n", chain)
	if err != nil {
		t.Errorf("Error querying Delegator Delegations: %v", err)
	}
	require.Equal(t, reflect.TypeOf(types.DelegationResponses{}), reflect.TypeOf(queryResponse))
}

func TestQueryDelegatorValidatorDelegations(t *testing.T) {
	chain, _ := InitializeAndStartChain("336h", constants.DefaultPBridgeHome)
	query, err := QueryDelegatorValidatorDelegations("cosmos1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d7lzl74n", "cosmosvaloper1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d76ktteq", chain)
	if err != nil {
		t.Errorf("Error quering delegator validation delegations: %v", err)
	}
	reDelegator := regexp.MustCompile(`^cosmos[0-9a-fA-F]`)
	reValidator := regexp.MustCompile(`^cosmosvaloper[0-9a-fA-F]`)
	require.Equal(t, true, reDelegator.MatchString(query.GetDelegation().DelegatorAddress))
	require.Equal(t, true, reValidator.MatchString(query.GetDelegation().ValidatorAddress))
	require.NotNil(t, query)
	require.Equal(t, reflect.TypeOf(types.DelegationResponse{}), reflect.TypeOf(query))
}

func TestQueryValidatorDelegator(t *testing.T) {
	chain, _ := InitializeAndStartChain("336h", constants.DefaultPBridgeHome)
	query, err := QueryValidatorDelegator("cosmos1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d7lzl74n", "cosmosvaloper1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d76ktteq", chain)
	if err != nil {
		t.Errorf("Error Quering Validator Delegator: %v", err)
	}
	re := regexp.MustCompile(`^cosmosvaloper[0-9a-fA-F]`)
	require.Equal(t, true, re.MatchString(query.OperatorAddress))
	require.Equal(t, reflect.TypeOf(types.Validator{}), reflect.TypeOf(query))
}
