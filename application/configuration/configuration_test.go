//go:build units

/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package configuration

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
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

// fixme tests depends on each other
func TestValidateAndSeal(t *testing.T) {
	SetConfig(test.GetCmdWithConfig())

	wrapAddress, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	SetCASPAddresses(wrapAddress, common.HexToAddress("0x5988ab40c82bbbb2067eec1e19b08cdc8d5e22d5"))

	ValidateAndSeal()
	require.Equal(t, GetAppConfig().seal, true, "appConfig did not get validated")
}

func TestGetWrapAddress(t *testing.T) {
	config := SetConfig(test.GetCmdWithConfig())
	wrapAddress, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	require.Equal(t, config.Tendermint.GetWrapAddress(), wrapAddress.String(), "wrapAddress not set")
}
