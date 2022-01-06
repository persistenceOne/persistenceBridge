/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package casp

import (
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestGetEthAddress(t *testing.T) {
	test.SetTestConfig()
	ethAddress, err := GetEthAddress()
	re := regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
	require.Nil(t, err)
	require.NotNil(t, ethAddress)
	require.Equal(t, 20, len(ethAddress))
	require.Equal(t, true, re.MatchString(ethAddress.String()))
}

func TestGetTendermintAddress(t *testing.T) {
	test.SetTestConfig()
	tenderMintAddress, errTMA := GetTendermintAddress()
	re := regexp.MustCompile(`^cosmos[0-9a-zA-Z]{39}$`)
	require.Nil(t, errTMA, "Error Getting Tendermint address")
	require.Equal(t, true, re.MatchString(tenderMintAddress.String()))
	require.NotNil(t, tenderMintAddress)
	require.Equal(t, 20, len(tenderMintAddress))
}
