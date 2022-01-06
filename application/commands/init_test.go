/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/stretchr/testify/require"
)

func TestInitCommand(t *testing.T) {
	cmd := InitCommand()
	err := cmd.Flags().Set(constants.FlagPBridgeHome, constants.TestHomeDir)
	require.Equal(t, nil, err)
	err = cmd.Execute()
	require.Equal(t, nil, err)

	config := configuration.GetAppConfig()
	_, err = toml.DecodeFile(filepath.Join(constants.TestHomeDir, "config.toml"), &config)
	require.Equal(t, nil, err)

}
