/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package casp

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/utilities/test"
)

func TestGetUncompressedPublicKeys(t *testing.T) {
	configuration.InitConfig()
	configuration.SetConfig(test.GetCmdWithConfig())

	funcResponse, err := GetUncompressedTMPublicKeys()
	require.Nil(t, err)
	require.Equal(t, funcResponse.AccountName, "tendermint")
	require.Equal(t, 1, len(funcResponse.Items))

	funcResponse, err = GetUncompressedEthPublicKeys()
	require.Nil(t, err)
	require.Equal(t, funcResponse.AccountName, "ethereum")
	require.Equal(t, 1, len(funcResponse.Items))
}
