/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package testing

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/utilities/testing/cmd"
	"io/ioutil"
	"os"
	"path/filepath"
)

func SetTestConfig() {
	config := configuration.SetConfig(cmd.GetCmdWithConfig())
	configPath := filepath.Join(constants2.TestHomeDir, "config.toml")
	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	if err := encoder.Encode(config); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(constants2.TestHomeDir, os.ModePerm); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(configPath, buf.Bytes(), 0644); err != nil {
		panic(err)
	}
}
