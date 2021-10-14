package tendermint

import (
	"github.com/persistenceOne/persistenceBridge/application/constants"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddressIsDelegatorToValidator(t *testing.T) {
	test.SetTestConfig()
	chain, _ := InitializeAndStartChain("336h", constants.DefaultPBridgeHome)
	delegatorToValidator := AddressIsDelegatorToValidator("cosmos1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d7lzl74n", "cosmosvaloper1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d76ktteq", chain)
	require.Equal(t, false, delegatorToValidator)
}

func TestQueryDelegatorDelegations(t *testing.T) {
	test.SetTestConfig()
	chain, _ := InitializeAndStartChain("336h", constants.DefaultPBridgeHome)
	_, err := QueryDelegatorDelegations("cosmos1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d7lzl74n", chain)
	require.Equal(t, "rpc error: code = NotFound desc = rpc error: code = NotFound desc = unable to find delegations for address cosmos1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d7lzl74n: key not found", err.Error())
}

func TestQueryDelegatorValidatorDelegations(t *testing.T) {
	test.SetTestConfig()
	chain, _ := InitializeAndStartChain("336h", constants.DefaultPBridgeHome)
	_, err := QueryDelegatorValidatorDelegations("cosmos1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d7lzl74n", "cosmosvaloper1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d76ktteq", chain)
	require.Equal(t, "rpc error: code = NotFound desc = rpc error: code = NotFound desc = delegation with delegator cosmos1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d7lzl74n not found for validator cosmosvaloper1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d76ktteq: key not found", err.Error())

}

func TestQueryValidatorDelegator(t *testing.T) {
	test.SetTestConfig()
	chain, _ := InitializeAndStartChain("336h", constants.DefaultPBridgeHome)
	_, err := QueryValidatorDelegator("cosmos1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d7lzl74n", "cosmosvaloper1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d76ktteq", chain)
	require.Equal(t, "rpc error: code = Unknown desc = rpc error: code = Internal desc = no delegation for (address, validator) tuple: unknown request", err.Error())
}
