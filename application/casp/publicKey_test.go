package casp

import (
	"crypto/ecdsa"
	"github.com/BurntSushi/toml"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	caspQueries "github.com/persistenceOne/persistenceBridge/application/rest/casp"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"
)

func TestGetTMPubKey(t *testing.T) {
	configuration.InitConfig()
	appConfig := test.GetCmdWithConfig()
	configuration.SetConfig(&appConfig)
	
	uncompressedPublicKeys, err := caspQueries.GetUncompressedEthPublicKeys()
	if err != nil {
		t.Errorf("Failed to get casp Response")
	}
	tmpKey := GetTMPubKey(uncompressedPublicKeys.PublicKeys[0])
	re := regexp.MustCompile(`^PubKeySecp256k1{+[0-9a-fA-F]+}$`)
	require.Equal(t, 20, len(tmpKey.Address().Bytes()))
	require.Equal(t, reflect.TypeOf(types.Address{}),reflect.TypeOf(tmpKey.Address()))
	require.Equal(t, true,re.MatchString(tmpKey.String()),"TM Public Key regex not matching")
	require.NotNil(t, tmpKey)
}

func TestGetEthPubKey(t *testing.T) {
	configuration.InitConfig()
	appConfig := test.GetCmdWithConfig()
	configuration.SetConfig(&appConfig)
	
	uncompressedPublicKeys, err := caspQueries.GetUncompressedEthPublicKeys()
	if err != nil {
		t.Errorf("Failed to get casp Response")
	}
	ethPubliKey := uncompressedPublicKeys.PublicKeys[0]
	ethKey := GetEthPubKey(ethPubliKey)
	require.Equal(t, 20, len(crypto.PubkeyToAddress(ethKey)))
	require.Equal(t, reflect.TypeOf(ecdsa.PublicKey{}),reflect.TypeOf(ethKey))
	require.Equal(t, reflect.TypeOf(&big.Int{}),reflect.TypeOf(ethKey.X))
	require.Equal(t, reflect.TypeOf(&big.Int{}),reflect.TypeOf(ethKey.Y))
	require.NotNil(t, ethKey)
}

func TestGetXY(t *testing.T)  {
	configuration.InitConfig()
	appConfig := test.GetCmdWithConfig()
	configuration.SetConfig(&appConfig)
	
	uncompressedPublicKeys, err := caspQueries.GetUncompressedEthPublicKeys()
	if err != nil {
		t.Errorf("Failed to get casp Response")
	}
	x, y := getXY(uncompressedPublicKeys.PublicKeys[0])
	require.Equal(t, 32, len(y.Bytes()))
	require.Equal(t, 32, len(y.Bytes()))
	require.Equal(t, reflect.TypeOf(big.Int{}),x)
	require.Equal(t, reflect.TypeOf(big.Int{}),y)
	require.NotNil(t, x)
	require.NotNil(t, y)
}
