package db

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestInitializeDB(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	dbPath := (filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")

	var ethStart int64 = 4772131
	var tmStart int64 = 1
	_, err := InitializeDB(dbPath, tmStart, ethStart)
	require.Nil(t, err)
	ethStatus, err := GetEthereumStatus()
	require.Nil(t, err)
	cosmosLastCheckHeight, err := GetCosmosStatus()
	require.Nil(t, err)
	ethHeight := ethStatus.LastCheckHeight + 1
	cosmosHeight := cosmosLastCheckHeight.LastCheckHeight + 1
	require.Equal(t, ethStart, ethHeight)
	require.Equal(t, tmStart, cosmosHeight )

	db.Close()

}

func TestOpenDB(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	require.Nil(t, err)

	db.Close()
}
