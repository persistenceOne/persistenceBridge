/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package configuration

import (
	"encoding/json"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceBridge/utilities/test"
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

func TestSetConfigAndChange(t *testing.T) {
	InitConfig()

	appConfig := SetConfig(test.GetCmdWithConfig())
	require.Equal(t, GetAppConfig(), *appConfig)

	oldConfigBytes, err := json.Marshal(appConfig)
	require.Nil(t, err)

	GetAppConfig().Kafka.TopicDetail.ReplicaAssignment = map[int32][]int32{
		99: {100},
	}

	oldConfig := new(config)

	err = json.Unmarshal(oldConfigBytes, oldConfig)
	require.Nil(t, err)

	require.Equal(t, appConfig, oldConfig)
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
