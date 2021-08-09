package casp

import (
	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"path/filepath"
	"testing"
)

var ethBridgeAdmin common.Address

func TestGetSignOperation(t *testing.T) {
	dirname, err := os.UserHomeDir()

	pStakeConfig := configuration.InitConfig()
	_, err = toml.DecodeFile(filepath.Join(dirname, "/.persistenceBridge/config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}

	operationID := "69544933-2767-4e09-af4a-c2dacb9a20af"
	reponseRecieved, err := GetSignOperation(operationID)

	require.Equal(t, nil, err)
	require.Equal(t, reponseRecieved.StatusText, "Completed")
	require.Equal(t, reponseRecieved.IsApproved, true)
	require.Equal(t, reponseRecieved.AccountID, "bd4c618e-8046-4fef-bdaa-9716ade77553")

}
