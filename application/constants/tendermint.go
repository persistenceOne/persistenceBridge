package constants

import sdkTypes "github.com/cosmos/cosmos-sdk/types"

const (
	validator1 = "cosmosvaloper1pkkayn066msg6kn33wnl5srhdt3tnu2v8fyhft"

	TxEvents    = "tm.event='Tx'"
	BlockEvents = "tm.event='NewBlock'"
	BlockHeader = "tm.event='NewBlockHeader'"
)

var (
	Validator1, _ = sdkTypes.ValAddressFromBech32(validator1)
	MinimumAmount = sdkTypes.NewInt(int64(2000000))
)
