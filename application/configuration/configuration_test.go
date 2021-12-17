/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package configuration

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetAppConfig(t *testing.T) {
	InitConfig()
	require.Equal(t, *appConfig, GetAppConfig(), "The two configurations should be the same")
}

func TestInitConfig(t *testing.T) {
	newAppConfig := InitConfig()
	appConfig := newConfig()

	require.Equal(t, appConfig, *newAppConfig)
}

func TestSetConfig(t *testing.T) {
	InitConfig()
	appConfig := SetConfig(test.GetCmdWithConfig())
	require.Equal(t, GetAppConfig(), *appConfig)
}

func TestSetPStakeAddress(t *testing.T) {
	InitConfig()
	config := SetConfig(test.GetCmdWithConfig())
	pStakeAddress, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	SetPStakeAddress(pStakeAddress)
	require.Equal(t, config.Tendermint.pStakeAddress, pStakeAddress.String(), "PStakeAddress not set")
}

func TestValidateAndSeal(t *testing.T) {
	InitConfig()
	config := SetConfig(test.GetCmdWithConfig())
	pStakeAddress, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	SetPStakeAddress(pStakeAddress)
	ValidateAndSeal()
	require.Equal(t, config.seal, true, "appConfig did not get validated")
}

func TestGetPStakeAddress(t *testing.T) {
	InitConfig()
	config := SetConfig(test.GetCmdWithConfig())
	pStakeAddress, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	SetPStakeAddress(pStakeAddress)
	require.Equal(t, config.Tendermint.GetPStakeAddress(), pStakeAddress.String(), "pStakeAddress not set")
}
