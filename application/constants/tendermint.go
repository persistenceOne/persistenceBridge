package constants

import sdkTypes "github.com/cosmos/cosmos-sdk/types"

const (
	TxEvents    = "tm.event='Tx'"
	BlockEvents = "tm.event='NewBlock'"
	BlockHeader = "tm.event='NewBlockHeader'"
)

var (
	MinimumAmount = sdkTypes.NewInt(int64(5000000))
)
