package casp

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/stretchr/testify/require"
	"log"
	"path/filepath"
	"testing"
)

func TestGetEthAddress(t *testing.T) {
	pStakeConfig := configuration.InitConfig()
	_, err := toml.DecodeFile(filepath.Join("/Users/ankitkumar/.persistenceBridge/", "config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}
	ethAddress, err := GetEthAddress()
	if err != nil {
		log.Fatalf("An error Occured %v",err)
	}
	fmt.Println(ethAddress)
	require.NotNil(t,ethAddress)
	fmt.Println(ethAddress.Bytes())
	require.Equal(t,20,len(ethAddress))
}

func TestGetTendermintAddress(t *testing.T) {
	pStakeConfig := configuration.InitConfig()
	_, err := toml.DecodeFile(filepath.Join("/Users/ankitkumar/.persistenceBridge/", "config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}
	tenderMintAddress, error := GetTendermintAddress()
	if error != nil {
		t.Errorf("Failed to get tindermint Address")
	}
	fmt.Println(tenderMintAddress)
	require.NotNil(t, tenderMintAddress)
	require.Equal(t, 20, len(tenderMintAddress))
}
