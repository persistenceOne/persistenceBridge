package constants

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	TestHomeDir = os.ExpandEnv("$HOME/testPersistenceBridge")
	TestDbDir   = os.ExpandEnv("$HOME/testPersistenceBridge/db")
)

func LoadEnv() {
	home, err := os.UserHomeDir()
	envDir := home + "/testPersistenceBridge/.env"
	err = godotenv.Load(envDir)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
