//go:build integration

/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package casp

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/utilities/test"
)

func TestGetEthAddress(t *testing.T) {
	configuration.SetConfig(test.GetCmdWithConfig())

	ethAddress, err := GetEthAddress()
	require.Nil(t, err)
	require.NotNil(t, ethAddress)
	require.Equal(t, 20, len(ethAddress))

	re := regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
	require.Equal(t, true, re.MatchString(ethAddress.String()))
}

func TestGetTendermintAddress(t *testing.T) {
	configuration.SetConfig(test.GetCmdWithConfig())

	tenderMintAddress, errTMA := GetTendermintAddress()
	require.Nil(t, errTMA, "Error Getting Tendermint address")
	require.NotNil(t, tenderMintAddress)
	require.Equal(t, 20, len(tenderMintAddress))

	re := regexp.MustCompile(`^cosmos[0-9a-zA-Z]{39}$`)
	require.Equal(t, true, re.MatchString(tenderMintAddress.String()))
}
