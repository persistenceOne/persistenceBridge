package casp

import (
	"github.com/BurntSushi/toml"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func TestGetEthAddress(t *testing.T) {
	pStakeConfig := configuration.InitConfig()
	dirname, _ := os.UserHomeDir()
	_, err := toml.DecodeFile(filepath.Join(dirname, "/.persistenceBridge/config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}
	ethAddress, err := GetEthAddress()
	re := regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
	require.Nil(t, err)
	require.NotNil(t,ethAddress)
	require.Equal(t,20,len(ethAddress))
	require.Equal(t,true,re.MatchString(ethAddress.String()) )
}

func TestGetTendermintAddress(t *testing.T) {
	pStakeConfig := configuration.InitConfig()
	dirname, _ := os.UserHomeDir()
	_, err := toml.DecodeFile(filepath.Join(dirname, "/.persistenceBridge/config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}
	tenderMintAddress, errTMA := GetTendermintAddress()
	re := regexp.MustCompile(`^cosmos[0-9a-zA-Z]{39}$`)
	require.Nil(t, errTMA,"Error Getting Tendermint address")
	require.Equal(t, true,re.MatchString(tenderMintAddress.String()))
	require.NotNil(t, tenderMintAddress)
	require.Equal(t, 20, len(tenderMintAddress))
}
