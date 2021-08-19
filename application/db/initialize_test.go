package db

import (
	"github.com/BurntSushi/toml"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestInitializeDB(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	pStakeConfig := configuration.InitConfig()
	_, err := toml.DecodeFile(filepath.Join(dirname, "/.persistenceBridge/config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}

	dbPath := (filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")

	var ethStart int64 = 4772131
	var tmStart int64 = 1
	_, err = InitializeDB(dbPath, tmStart, ethStart)
	if err != nil {
		t.Fatalf("Error %v", err.Error())
	}
}

func TestOpenDB(t *testing.T) {
	db.Close()
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	db.Close()
}
