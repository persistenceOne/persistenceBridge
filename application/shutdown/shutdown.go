/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package shutdown

var seal = false
var stopBridge = false
var tendermintStopped = false
var ethStopped = false
var kafkaConsumerClosed = false

func GetBridgeStopSignal() bool {
	return stopBridge
}

func StopBridge() {
	if !seal {
		stopBridge = true
		seal = true
	}
}

func GetTMStopped() bool {
	return tendermintStopped
}

func SetTMStopped(value bool) {
	tendermintStopped = value
}

func GetETHStopped() bool {
	return ethStopped
}

func SetETHStopped(value bool) {
	ethStopped = value
}

func GetKafkaConsumerClosed() bool {
	return kafkaConsumerClosed
}

func SetKafkaConsumerClosed(value bool) {
	kafkaConsumerClosed = value
}
