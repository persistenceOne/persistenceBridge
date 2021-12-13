/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package constants

import "os"

var (
	TestHomeDir = os.ExpandEnv("$HOME/testPersistenceBridge")
	TestDbDir   = os.ExpandEnv("$HOME/testPersistenceBridge/db")
)
