package contracts

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)



func TestOnStake(t *testing.T) {
	configuration.InitConfig()
	appConfig := test.GetCmdWithConfig()
	configuration.SetConfig(&appConfig)
	tmAddress, err := casp.GetTendermintAddress()
	require.Equal(t, nil, err)

	configuration.SetPStakeAddress(tmAddress)

	i := new(big.Int)
	i.SetInt64(1000)
	arr := []interface{}{common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa")),  i}
	stakeMsg, ercAddress, err := onStake(arr)
	require.Equal(t, nil, err)
	require.Equal(t, common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa")).String(), ercAddress.String())
	stakeMsgString := stakeMsg.String()
	require.NotNil(t, stakeMsgString)
}

func TestOnUnStake(t *testing.T) {
	configuration.InitConfig()
	appConfig := test.GetCmdWithConfig()
	configuration.SetConfig(&appConfig)
	tmAddress, err := casp.GetTendermintAddress()
	require.Equal(t, nil, err)

	configuration.SetPStakeAddress(tmAddress)

	i := new(big.Int)
	i.SetInt64(1000)
	arr := []interface{}{common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa")),  i}
	UnStakeMsg, ercAddress, err := onUnStake(arr)
	require.Equal(t, nil, err)
	require.Equal(t, common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa")).String(), ercAddress.String())
	UnStakeMsgString := UnStakeMsg.String()
	require.NotNil(t, UnStakeMsgString)
}
