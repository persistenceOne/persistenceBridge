/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package casp

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	caspQueries "github.com/persistenceOne/persistenceBridge/application/rest/casp"
)

func GetEthAddress(ctx context.Context) (common.Address, error) {
	uncompressedPublicKeys, err := caspQueries.GetUncompressedEthPublicKeys(ctx)
	if err != nil {
		return common.Address{}, err
	}

	if len(uncompressedPublicKeys.Items) == 0 {
		return common.Address{}, fmt.Errorf("%w: ethereum", ErrNoPublicKeys)
	}

	publicKey := GetEthPubKey(uncompressedPublicKeys.Items[0])

	return crypto.PubkeyToAddress(publicKey), nil
}

func GetTendermintAddress(ctx context.Context) (sdk.AccAddress, error) {
	uncompressedPublicKeys, err := caspQueries.GetUncompressedTMPublicKeys(ctx)
	if err != nil {
		return nil, err
	}

	if len(uncompressedPublicKeys.Items) == 0 {
		return nil, fmt.Errorf("%w: tendermint", ErrNoPublicKeys)
	}

	tmPublicKey := GetTMPubKey(uncompressedPublicKeys.Items[0])

	return sdk.AccAddress(tmPublicKey.Address()), nil
}
