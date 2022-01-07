/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package utils

// Consumer groups

const (
	GroupToEth           = "group-to-ethereum"
	GroupToTendermint    = "group-to-tendermint"
	GroupEthUnbond       = "group-ethereum-unbond"
	GroupMsgSend         = "group-msg-send"
	GroupMsgDelegate     = "group-msg-delegate"
	GroupRedelegate      = "group-redelegate"
	GroupMsgUnbond       = "group-msg-unbond"
	GroupRetryTendermint = "group-retry-tendermint"
)

func Groups() []string {
	return []string{GroupEthUnbond,
		GroupMsgSend, GroupMsgDelegate, GroupRedelegate, GroupMsgUnbond,
		GroupToEth, GroupRetryTendermint, GroupToTendermint,
	}
}

// Topics
const (
	ToEth           = "to-ethereum"
	ToTendermint    = "to-tendermint"
	MsgSend         = "msg-send"        // priority3
	MsgDelegate     = "msg-delegate"    // priority2
	MsgUnbond       = "msg-unbond"      // priority1
	EthUnbond       = "ethereum-unbond" // flushes every 3 days
	Redelegate      = "redelegate"
	RetryTendermint = "retry-tendermint"
)

// Topics : is list of topics
func Topics() []string {
	return []string{
		EthUnbond, MsgSend, MsgDelegate, Redelegate, MsgUnbond,
		ToEth, RetryTendermint, ToTendermint,
	}
}
