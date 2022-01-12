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

func TestOnWithdrawUTokens(t *testing.T) {
	configuration.SetConfig(test.GetCmdWithConfig())

	ctx := context.Background()

	tmAddress, err := casp.GetTendermintAddress(ctx)
	require.Nil(t, err)

	ethAddress, err := casp.GetEthAddress(ctx)
	require.Nil(t, err)

	configuration.SetCASPAddresses(tmAddress, ethAddress)

	i := new(big.Int)
	i.SetInt64(1000)
	arr := []interface{}{common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa")), i, "cosmos1aqxce9xssqsyjfm3gt39w4gf9u9dxgax6qjk79"}

	sendCoinMsg, ercAddress, err := onWithdrawUTokens(arr)
	require.Nil(t, err)
	require.Equal(t, common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa")).String(), ercAddress.String())

	sendCoinMsgString := sendCoinMsg.String()
	require.NotNil(t, sendCoinMsgString)

	arr = []interface{}{common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa")), i, ""}

	sendCoinMsg, ercAddress, err = onWithdrawUTokens(arr)
	require.Equal(t, "empty address string is not allowed", err.Error())
}
