/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package testing

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/utilities/testing/cmd"
)

func SetTestConfig() {
	config := configuration.SetConfig(cmd.GetCmdWithConfig())
	configPath := filepath.Join(constants.TestHomeDir, "config.toml")
	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	if err := encoder.Encode(config); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(constants.TestHomeDir, os.ModePerm); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(configPath, buf.Bytes(), 0644); err != nil {
		panic(err)
	}
}
