package constants

import sdkTypes "github.com/cosmos/cosmos-sdk/types"

const (
	validator1 = "cosmosvaloper19ylp572rf3f8lc9hl2pjv5nfckmgj3sj8kkqma"

	TxEvents    = "tm.event='Tx'"
	BlockEvents = "tm.event='NewBlock'"
	BlockHeader = "tm.event='NewBlockHeader'"
)

var (
	Validator1, _ = sdkTypes.ValAddressFromBech32(validator1)
	MinimumAmount = sdkTypes.NewInt(int64(5000000))
)
