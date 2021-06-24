package casp

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"log"
	"math/big"
)

// Should include prefix "04"
func GetPubKey(caspPubKey string) cryptotypes.PubKey {

	pubKeyBytes, err := hex.DecodeString(string([]rune(caspPubKey)[2:])) // uncompressed pubkey
	if err != nil {
		log.Fatalln(err)
	}
	var x big.Int
	x.SetBytes(pubKeyBytes[0:32])
	var y big.Int
	y.SetBytes(pubKeyBytes[32:])

	pubKey := ecdsa.PublicKey{
		Curve: btcec.S256(),
		X:     &x,
		Y:     &y,
	}
	pubkeyObject := (*btcec.PublicKey)(&pubKey)
	pk := pubkeyObject.SerializeCompressed()
	return &secp256k1.PubKey{Key: pk}
}
