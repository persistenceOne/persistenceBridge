/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package constants

import (
	"os"
)

// nolint safe private global vars with access by public function
// nolint: gochecknoglobals
var (
	testHomeDir = os.ExpandEnv("$HOME/testPersistenceBridge")
)

func TestHomeDir() string {
	return testHomeDir
}
