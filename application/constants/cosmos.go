package constants

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var (
	MsgSendTypeUrl                    = sdk.MsgTypeURL(&bankTypes.MsgSend{})
	MsgDelegateTypeUrl                = sdk.MsgTypeURL(&stakingTypes.MsgDelegate{})
	MsgUndelegateTypeUrl              = sdk.MsgTypeURL(&stakingTypes.MsgUndelegate{})
	MsgWithdrawDelegatorRewardTypeUrl = sdk.MsgTypeURL(&distributionTypes.MsgWithdrawDelegatorReward{})
)
