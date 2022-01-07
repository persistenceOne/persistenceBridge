/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package casp

import (
	"github.com/ethereum/go-ethereum/crypto"
	caspQueries "github.com/persistenceOne/persistenceBridge/application/rest/casp"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestGetTMPubKey(t *testing.T) {
	test.SetTestConfig()

	uncompressedPublicKeys, err := caspQueries.GetUncompressedTMPublicKeys()
	require.Nil(t, err, "Failed to get casp Response")
	tmpKey := GetTMPubKey(uncompressedPublicKeys.Items[0])
	re := regexp.MustCompile(`^PubKeySecp256k1{+[0-9a-fA-F]+}$`)
	require.Equal(t, 20, len(tmpKey.Address().Bytes()))
	require.Equal(t, true, re.MatchString(tmpKey.String()), "TM Public Key regex not matching")
	require.NotNil(t, tmpKey)
}

func TestGetEthPubKey(t *testing.T) {
	test.SetTestConfig()

	uncompressedPublicKeys, err := caspQueries.GetUncompressedEthPublicKeys()
	require.Nil(t, err, "Failed to get casp Response")
	ethPublicKey := GetEthPubKey(uncompressedPublicKeys.Items[0])
	require.Equal(t, 20, len(crypto.PubkeyToAddress(ethPublicKey)))
	require.NotNil(t, ethPublicKey)
}

func TestGetXY(t *testing.T) {
	test.SetTestConfig()

	uncompressedPublicKeys, err := caspQueries.GetUncompressedEthPublicKeys()
	require.Nil(t, err, "Failed to get casp Response")
	x, y := getXY(uncompressedPublicKeys.Items[0])
	require.NotNil(t, x)
	require.NotNil(t, y)
	require.LessOrEqual(t, 32, len(x.Bytes()))
	require.LessOrEqual(t, 32, len(y.Bytes()))
}
