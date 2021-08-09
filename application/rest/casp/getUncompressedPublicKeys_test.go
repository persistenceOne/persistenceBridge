package casp

import (
	"github.com/BurntSushi/toml"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"path/filepath"
	"testing"
)


func Test_getUncompressedPublicKeys(t *testing.T) {
	dirname, err := os.UserHomeDir()

	pStakeConfig := configuration.InitConfig()
	_, err = toml.DecodeFile(filepath.Join(dirname, "/.persistenceBridge/config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}

	funcResponse, err := getUncompressedPublicKeys(118)
	require.Equal(t, nil, err)
	require.Equal(t, funcResponse.AccountName, "tendermint")
	require.Equal(t, 1, len(funcResponse.PublicKeys))


	funcResponse, err = getUncompressedPublicKeys(60)
	require.Equal(t, nil, err)
	require.Equal(t, funcResponse.AccountName, "ethereum")

}
