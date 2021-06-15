package constants

import sdkTypes "github.com/cosmos/cosmos-sdk/types"

const (
	validator1 = "cosmosvaloper1phzsnju4t0alc9r8rc9jfyzpwf22jsaup9dz4d"

	TxEvents    = "tm.event='Tx'"
	BlockEvents = "tm.event='NewBlock'"
	BlockHeader = "tm.event='NewBlockHeader'"
)

var (
	Validator1, _ = sdkTypes.ValAddressFromBech32(validator1)
	MinimumAmount = sdkTypes.NewInt(int64(2000000))
)
