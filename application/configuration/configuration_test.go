/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package configuration

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceBridge/utilities/test"
)

func TestGetAppConfig(t *testing.T) {
	require.Equal(t, appConfig, GetAppConfig(), "The two configurations should be the same")
}

func TestSetConfig(t *testing.T) {
	appConfig := SetConfig(test.GetCmdWithConfig())
	require.Equal(t, GetAppConfig(), appConfig)
}

func TestSetConfigAndChange(t *testing.T) {
	appConfig := SetConfig(test.GetCmdWithConfig())
	require.Equal(t, GetAppConfig(), appConfig)

	appConfigOld := appConfig.DeepCopy()

	newConfig := GetAppConfig()
	newConfig.Kafka.Brokers[0] = "100000"

	require.Equal(t, appConfig, appConfigOld)
	require.NotEqual(t, appConfig, newConfig)
}

func TestSetPStakeAddress(t *testing.T) {
	config := SetConfig(test.GetCmdWithConfig())
	pStakeAddress, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")

	SetPStakeAddress(pStakeAddress)
	require.Equal(t, config.Tendermint.pStakeAddress, pStakeAddress.String(), "PStakeAddress not set")
}

func TestValidateAndSeal(t *testing.T) {
	config := SetConfig(test.GetCmdWithConfig())
	pStakeAddress, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")

	SetPStakeAddress(pStakeAddress)
	ValidateAndSeal()
	require.Equal(t, config.seal, true, "appConfig did not get validated")
}

func TestGetPStakeAddress(t *testing.T) {
	config := SetConfig(test.GetCmdWithConfig())

	pStakeAddress, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	SetPStakeAddress(pStakeAddress)
	require.Equal(t, config.Tendermint.GetPStakeAddress(), pStakeAddress.String(), "pStakeAddress not set")
}
