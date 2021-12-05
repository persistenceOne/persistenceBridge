package constants

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

var (
	MsgSendTypeUrl                    = sdk.MsgTypeURL(&bankTypes.MsgSend{})
	MsgWithdrawDelegatorRewardTypeUrl = sdk.MsgTypeURL(&distributionTypes.MsgWithdrawDelegatorReward{})
)
