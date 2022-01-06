/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

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

func TestOnWithdrawUTokens(t *testing.T) {
	test.SetTestConfig()
	tmAddress, err := casp.GetTendermintAddress()
	require.Equal(t, nil, err)
	ethAddress, err := casp.GetEthAddress()
	require.Equal(t, nil, err)
	configuration.SetCASPAddresses(tmAddress, ethAddress)

	i := new(big.Int)
	i.SetInt64(1000)
	arr := []interface{}{common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa")), i, "cosmos1aqxce9xssqsyjfm3gt39w4gf9u9dxgax6qjk79"}
	sendCoinMsg, ercAddress, err := onWithdrawUTokens(arr)
	require.Equal(t, nil, err)
	require.Equal(t, common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa")).String(), ercAddress.String())
	sendCoinMsgString := sendCoinMsg.String()
	require.NotNil(t, sendCoinMsgString)

	arr = []interface{}{common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa")), i, ""}
	sendCoinMsg, ercAddress, err = onWithdrawUTokens(arr)
	require.Equal(t, "empty address string is not allowed", err.Error())
}
