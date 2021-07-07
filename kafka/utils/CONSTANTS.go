/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package utils

// Consumer groups

const GroupToEth = "group-to-ethereum"
const GroupToTendermint = "group-to-tendermint"
const GroupEthUnbond = "group-ethereum-unbond"
const GroupMsgSend = "group-msg-send"
const GroupMsgDelegate = "group-msg-delegate"
const GroupMsgUnbond = "group-msg-unbond"

var Groups = []string{GroupEthUnbond,
	GroupMsgSend, GroupMsgDelegate, GroupMsgUnbond,
	GroupToEth, GroupToTendermint,
}

//Topics

const ToEth = "to-ethereum"
const ToTendermint = "to-tendermint"
const MsgSend = "msg-send"          //priority3
const MsgDelegate = "msg-delegate"  //priority2
const MsgUnbond = "msg-unbond"      //priority1
const EthUnbond = "ethereum-unbond" //flushes every 3 days
const Redelegate = "redelegate"

// Topics : is list of topics
var Topics = []string{
	EthUnbond, MsgSend, MsgDelegate, MsgUnbond,
	ToEth, ToTendermint,
}
