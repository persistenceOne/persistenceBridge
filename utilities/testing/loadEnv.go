package testing

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadEnv() {
	home, err := os.UserHomeDir()
	envDir := home + "/testPersistenceBridge/.env"
	err = godotenv.Load(envDir)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
