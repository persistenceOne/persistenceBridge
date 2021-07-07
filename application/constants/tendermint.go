package constants

import sdkTypes "github.com/cosmos/cosmos-sdk/types"

const (
	validator1 = "cosmosvaloper1d20g0gcwhrwv8f2626dx0nkhauu0rsqkz4axrg"

	TxEvents    = "tm.event='Tx'"
	BlockEvents = "tm.event='NewBlock'"
	BlockHeader = "tm.event='NewBlockHeader'"
)

var (
	Validator1, _ = sdkTypes.ValAddressFromBech32(validator1)
	MinimumAmount = sdkTypes.NewInt(int64(5000000))
)
