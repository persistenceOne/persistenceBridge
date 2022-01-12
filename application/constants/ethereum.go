/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package constants

import (
	"github.com/ethereum/go-ethereum/common"
)

// nolint because those are not credentials and public
// nolint:gosec
const (
	LiquidStakingStake   = "stake"
	LiquidStakingUnStake = "unStake"

	TokenWrapperWithdrawUTokens = "withdrawUTokens"

	EthereumBlockConfirmations = 12
)


// nolint: gochecknoglobals
var (
	ethEmptyAddress = common.Address{}
	ethEmptyHash    = common.Hash{}
)

func EthereumZeroAddress() common.Address {
	return ethEmptyAddress
}

func EthereumEmptyTxHash() common.Hash {
	return ethEmptyHash
}