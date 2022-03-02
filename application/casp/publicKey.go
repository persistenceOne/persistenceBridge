/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package casp

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

// GetTMPubKey caspPubKey should include prefix "04"
func GetTMPubKey(caspPubKey string) cryptotypes.PubKey {
	x, y := getXY(caspPubKey)

	pubKey := ecdsa.PublicKey{
		Curve: btcec.S256(),
		X:     x,
		Y:     y,
	}
	if !(pubKey.Curve.IsOnCurve(pubKey.X, pubKey.Y)) {
		logging.Fatal("not a valid public key")
	}
	pubkeyObject := (*btcec.PublicKey)(&pubKey)
	pk := pubkeyObject.SerializeCompressed()
	return &secp256k1.PubKey{Key: pk}
}

// GetEthPubKey caspPubKey should include prefix "04"
func GetEthPubKey(caspPubKey string) ecdsa.PublicKey {
	x, y := getXY(caspPubKey)

	publicKey := ecdsa.PublicKey{
		Curve: crypto.S256(),
		X:     x,
		Y:     y,
	}

	if !(publicKey.Curve.IsOnCurve(publicKey.X, publicKey.Y)) {
		logging.Fatal("not a valid public key = ")
	}
	return publicKey

}

// getXY caspPubKey should include prefix "04"
func getXY(caspPubKey string) (x, y *big.Int) {
	if len(caspPubKey) < 2 {
		logging.Fatal("Invalid length of caspPubKey = " + caspPubKey)
	}
	s := strings.Split(caspPubKey, "")
	if s[0] != "0" && s[1] != "4" {
		logging.Fatal("invalid casp public key")
	}

	pubKeyBytes, err := hex.DecodeString(string([]rune(caspPubKey)[2:])) // uncompressed pubkey
	if err != nil {
		logging.Fatal(err)
	}

	if len(pubKeyBytes) != 64 {
		logging.Fatal(fmt.Sprintf("invalid casp public key, length (%v) not equal to 64", len(pubKeyBytes)))
	}

	x = big.NewInt(0).SetBytes(pubKeyBytes[0:32])
	y = big.NewInt(0).SetBytes(pubKeyBytes[32:])

	return
}
