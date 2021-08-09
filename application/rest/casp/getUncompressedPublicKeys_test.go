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


func Test_getUncompressedPublicKeys(t *testing.T) {
	pStakeConfig := configuration.InitConfig()
	_, err := toml.DecodeFile(filepath.Join("/Users/gokuls/.persistenceBridge/", "config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}


	funcResponse, err := getUncompressedPublicKeys(118)
	pubKey := configuration.GetAppConfig().CASP.EthereumPublicKey
	fmt.Println(pubKey)
	require.Equal(t, nil, err)
	require.Equal(t, funcResponse.AccountName, "tendermint")


	funcResponse, err = getUncompressedPublicKeys(60)
	require.Equal(t, nil, err)
	require.Equal(t, funcResponse.AccountName, "ethereum")
	//require.Equal(t, funcResponse.PublicKeys, )

}
