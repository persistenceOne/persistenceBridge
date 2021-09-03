package tendermint

import (
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
)


func TestAddressIsDelegatorToValidator(t *testing.T) {
	homedir, _ := os.UserHomeDir()
	fileName := strings.Join([]string{homedir,"/.persistenceBridge/chain.json"},"")
	chain, _ := InitializeAndStartChain(fileName, "336h", homedir)
	delegatorToValidator := AddressIsDelegatorToValidator("cosmos1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d7lzl74n", "cosmosvaloper1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d76ktteq", chain)
	require.Equal(t, true, delegatorToValidator)
}

func TestQueryDelegatorDelegations(t *testing.T) {
	homedir, _ := os.UserHomeDir()
	fileName := strings.Join([]string{homedir,"/.persistenceBridge/chain.json"},"")
	chain, _ := InitializeAndStartChain(fileName, "336h", homedir)
	queryResponse, err := QueryDelegatorDelegations("cosmos1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d7lzl74n",chain)
	if err != nil {
		t.Errorf("Error querying Delegator Delegations: %v",err)
	}
	require.Equal(t, reflect.TypeOf(types.DelegationResponses{}), reflect.TypeOf(queryResponse))
}

func TestQueryDelegatorValidatorDelegations(t *testing.T) {
	homedir, _ := os.UserHomeDir()
	fileName := strings.Join([]string{homedir,"/.persistenceBridge/chain.json"},"")
	chain, _ := InitializeAndStartChain(fileName, "336h", homedir)
	query, err := QueryDelegatorValidatorDelegations("cosmos1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d7lzl74n", "cosmosvaloper1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d76ktteq", chain)
	if err != nil {
		t.Errorf("Error quering delegator validation delegations: %v",err)
	}
	reDelegator := regexp.MustCompile(`^cosmos[0-9a-fA-F]`)
	reValidator := regexp.MustCompile(`^cosmosvaloper[0-9a-fA-F]`)
	require.Equal(t, true, reDelegator.MatchString(query.GetDelegation().DelegatorAddress))
	require.Equal(t, true, reValidator.MatchString(query.GetDelegation().ValidatorAddress))
	require.NotNil(t, query)
	require.Equal(t, reflect.TypeOf(types.DelegationResponse{}), reflect.TypeOf(query))
}

func TestQueryValidatorDelegator(t *testing.T) {
	homedir, _ := os.UserHomeDir()
	fileName := strings.Join([]string{homedir,"/.persistenceBridge/chain.json"},"")
	chain, _ := InitializeAndStartChain(fileName, "336h", homedir)
	query, err := QueryValidatorDelegator("cosmos1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d7lzl74n", "cosmosvaloper1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d76ktteq", chain)
	if err != nil {
		t.Errorf("Error Quering Validator Delegator: %v",err)
	}
	re := regexp.MustCompile(`^cosmosvaloper[0-9a-fA-F]`)
	require.Equal(t, true, re.MatchString(query.OperatorAddress))
	require.Equal(t, reflect.TypeOf(types.Validator{}), reflect.TypeOf(query))
}
