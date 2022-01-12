//go:build integration

/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package contracts

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/utilities/test"
)

func TestOnStake(t *testing.T) {
	configuration.SetConfig(test.GetCmdWithConfig())

	ctx := context.Background()

	tmAddress, err := casp.GetTendermintAddress(ctx)
	require.Nil(t, err)

	ethAddress, err := casp.GetEthAddress(ctx)
	require.Nil(t, err)

	configuration.SetCASPAddresses(tmAddress, ethAddress)

	i := new(big.Int)
	i.SetInt64(1000)
	arr := []interface{}{common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa")), i}
	stakeMsg, ercAddress, err := onStake(arr)
	require.Nil(t, err)

	require.Equal(t, common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa")).String(), ercAddress.String())

	stakeMsgString := stakeMsg.String()
	require.NotNil(t, stakeMsgString)
}

func TestOnUnStake(t *testing.T) {
	configuration.SetConfig(test.GetCmdWithConfig())

	ctx := context.Background()

	tmAddress, err := casp.GetTendermintAddress(ctx)
	require.Nil(t, err)

	ethAddress, err := casp.GetEthAddress(ctx)
	require.Nil(t, err)

	configuration.SetCASPAddresses(tmAddress, ethAddress)

	i := new(big.Int)
	i.SetInt64(1000)
	arr := []interface{}{common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa")), i}
	UnStakeMsg, ercAddress, err := onUnStake(arr)
	require.Nil(t, err)

	require.Equal(t, common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa")).String(), ercAddress.String())

	UnStakeMsgString := UnStakeMsg.String()
	require.NotNil(t, UnStakeMsgString)
}
