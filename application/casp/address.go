/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package casp

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	caspQueries "github.com/persistenceOne/persistenceBridge/application/rest/casp"
)

func GetEthAddress() (common.Address, error) {
	uncompressedPublicKeys, err := caspQueries.GetUncompressedEthPublicKeys()
	if err != nil {
		return common.Address{}, err
	}
	if len(uncompressedPublicKeys.Items) != 1 {
		return common.Address{}, fmt.Errorf("no or more than 1 eth public keys got from casp")
	}
	publicKey := GetEthPubKey(uncompressedPublicKeys.Items[0])
	return crypto.PubkeyToAddress(publicKey), nil
}

func GetTendermintAddress() (sdk.AccAddress, error) {
	uncompressedPublicKeys, err := caspQueries.GetUncompressedTMPublicKeys()
	if err != nil {
		return nil, err
	}
	if len(uncompressedPublicKeys.Items) != 1 {
		return nil, fmt.Errorf("no or more than 1 tendermint public keys got from casp")
	}
	tmPublicKey := GetTMPubKey(uncompressedPublicKeys.Items[0])
	return sdk.AccAddress(tmPublicKey.Address()), nil
}
