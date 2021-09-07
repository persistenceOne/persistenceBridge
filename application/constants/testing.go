package constants

import "os"

var (
	TestHomeDir = os.ExpandEnv("$HOME/testPersistenceBridge")
	TestDbDir   = os.ExpandEnv("$HOME/testPersistenceBridge/db")
)
