/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package casp

import (
	"crypto/ecdsa"
	"math/big"
	"reflect"
	"regexp"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	caspQueries "github.com/persistenceOne/persistenceBridge/application/rest/casp"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
)

func TestGetTMPubKey(t *testing.T) {
	configuration.InitConfig()
	configuration.SetConfig(test.GetCmdWithConfig())

	uncompressedPublicKeys, err := caspQueries.GetUncompressedEthPublicKeys()
	require.Nil(t, err, "Failed to get casp Response")

	tmpKey := GetTMPubKey(uncompressedPublicKeys.Items[0])
	re := regexp.MustCompile(`^PubKeySecp256k1{+[0-9a-fA-F]+}$`)

	require.Equal(t, 20, len(tmpKey.Address().Bytes()))
	require.Equal(t, reflect.TypeOf(types.Address{}), reflect.TypeOf(tmpKey.Address()))
	require.Equal(t, true, re.MatchString(tmpKey.String()), "TM Public Key regex not matching")
	require.NotNil(t, tmpKey)
}

func TestGetEthPubKey(t *testing.T) {
	configuration.InitConfig()
	configuration.SetConfig(test.GetCmdWithConfig())

	uncompressedPublicKeys, err := caspQueries.GetUncompressedEthPublicKeys()
	require.Nil(t, err, "Failed to get casp Response")

	ethPubliKey := uncompressedPublicKeys.Items[0]
	ethKey := GetEthPubKey(ethPubliKey)
	require.Equal(t, 20, len(crypto.PubkeyToAddress(ethKey)))
	require.Equal(t, reflect.TypeOf(ecdsa.PublicKey{}), reflect.TypeOf(ethKey))
	require.Equal(t, reflect.TypeOf(&big.Int{}), reflect.TypeOf(ethKey.X))
	require.Equal(t, reflect.TypeOf(&big.Int{}), reflect.TypeOf(ethKey.Y))
	require.NotNil(t, ethKey)
}

func TestGetXY(t *testing.T) {
	configuration.InitConfig()
	configuration.SetConfig(test.GetCmdWithConfig())

	uncompressedPublicKeys, err := caspQueries.GetUncompressedEthPublicKeys()
	require.Nil(t, err, "Failed to get casp Response")

	x, y := getXY(uncompressedPublicKeys.Items[0])
	require.NotNil(t, x)
	require.NotNil(t, y)
	require.Equal(t, 32, len(y.Bytes()))
	require.Equal(t, 32, len(y.Bytes()))
}
