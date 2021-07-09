package casp

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	caspQueries "github.com/persistenceOne/persistenceBridge/application/rest/casp"
)

func GetEthAddress() (common.Address, error) {
	uncompressedPublicKeys, err := caspQueries.GetUncompressedEthPublicKeys()
	if err != nil {
		return common.Address{}, err
	}
	if len(uncompressedPublicKeys.PublicKeys) == 0 {
		return common.Address{}, fmt.Errorf("no public keys got from casp")
	}
	publicKey := GetEthPubKey(uncompressedPublicKeys.PublicKeys[0])

	fromAddress := crypto.PubkeyToAddress(publicKey)
	return fromAddress, nil
}
