package constants

import sdkTypes "github.com/cosmos/cosmos-sdk/types"

const (
	validator1 = "cosmosvaloper1q8s35u88t5de6mlasaw9wugk5z2vd5y3ta4kc8"

	TxEvents    = "tm.event='Tx'"
	BlockEvents = "tm.event='NewBlock'"
	BlockHeader = "tm.event='NewBlockHeader'"
)

var (
	Validator1, _ = sdkTypes.ValAddressFromBech32(validator1)
	MinimumAmount = sdkTypes.NewInt(int64(5000000))
)
