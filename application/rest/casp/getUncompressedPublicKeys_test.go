//go:build integration

/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package casp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/utilities/test"
)

func TestGetUncompressedPublicKeys(t *testing.T) {
	configuration.SetConfig(test.GetCmdWithConfig())

	ctx := context.Background()

	funcResponse, err := GetUncompressedTMPublicKeys(ctx)
	require.Nil(t, err)
	require.Equal(t, funcResponse.AccountName, "tendermint")
	require.Equal(t, 1, len(funcResponse.Items))

	funcResponse, err = GetUncompressedEthPublicKeys(ctx)
	require.Nil(t, err)
	require.Equal(t, funcResponse.AccountName, "ethereum")
	require.Equal(t, 1, len(funcResponse.Items))
}
