/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package configuration

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	testCmd "github.com/persistenceOne/persistenceBridge/utilities/testing/cmd"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetAppConfig(t *testing.T) {
	require.Equal(t, appConfig, GetAppConfig(), "The two configurations should be the same")
}

func TestSetConfig(t *testing.T) {
	appConfig := SetConfig(testCmd.GetCmdWithConfig())
	require.Equal(t, GetAppConfig(), appConfig)
}

func TestValidateAndSeal(t *testing.T) {
	SetConfig(testCmd.GetCmdWithConfig())
	wrapAddress, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	SetCASPAddresses(wrapAddress, common.HexToAddress("0x5988ab40c82bbbb2067eec1e19b08cdc8d5e22d5"))
	ValidateAndSeal()
	require.Equal(t, appConfig.seal, true, "appConfig did not get validated")
}

func TestGetWrapAddress(t *testing.T) {
	SetConfig(testCmd.GetCmdWithConfig())
	pStakeAddress, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	require.Equal(t, appConfig.Tendermint.GetWrapAddress(), pStakeAddress.String(), "wrapAddress not set")
}
