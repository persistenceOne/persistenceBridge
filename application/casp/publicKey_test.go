package casp

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	caspQueries "github.com/persistenceOne/persistenceBridge/application/rest/casp"
	"github.com/stretchr/testify/require"
	"log"
	"path/filepath"
	"testing"
)

func TestGetTMPubKey(t *testing.T) {
	pStakeConfig := configuration.InitConfig()
	_, err := toml.DecodeFile(filepath.Join("/Users/ankitkumar/.persistenceBridge/", "config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}
	uncompressedPublicKeys, err := caspQueries.GetUncompressedEthPublicKeys()
	if err != nil {
		t.Errorf("Failed to get casp Response")
	}
	fmt.Println(uncompressedPublicKeys)
	tmpKey := GetTMPubKey(uncompressedPublicKeys.PublicKeys[0])
	fmt.Println(tmpKey.Address())
	require.Equal(t, 20, len(tmpKey.Address().Bytes()))
	require.NotNil(t, tmpKey)
}

func TestGetEthPubKey(t *testing.T) {
	pStakeConfig := configuration.InitConfig()
	_, err := toml.DecodeFile(filepath.Join("/Users/ankitkumar/.persistenceBridge/", "config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}
	uncompressedPublicKeys, err := caspQueries.GetUncompressedEthPublicKeys()
	if err != nil {
		t.Errorf("Failed to get casp Response")
	}
	ethKey := GetEthPubKey(uncompressedPublicKeys.PublicKeys[0])
	//fmt.Println(ethKey)
	require.Equal(t, 20, len(crypto.PubkeyToAddress(ethKey)))
	require.NotNil(t, ethKey)
}

func TestGetXY(t *testing.T)  {
	pStakeConfig := configuration.InitConfig()
	_, err := toml.DecodeFile(filepath.Join("/Users/ankitkumar/.persistenceBridge/", "config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}
	uncompressedPublicKeys, err := caspQueries.GetUncompressedEthPublicKeys()
	if err != nil {
		t.Errorf("Failed to get casp Response")
	}
	x, y := getXY(uncompressedPublicKeys.PublicKeys[0])
	fmt.Println(x.BitLen() , len(y.Bytes()))
	require.Equal(t, 32, len(y.Bytes()))
	require.Equal(t, 32, len(y.Bytes()))
	require.NotNil(t, x)
	require.NotNil(t, y)
	//fmt.Println(reflect.TypeOf(x))
	//require.Equal(t,big.Int{} ,x.)
}