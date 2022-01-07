/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package shutdown

// nolint fixme: each service must know its' status. remove global state
// nolint: gochecknoglobals
var (
	seal                = false
	stopBridge          = false
	tendermintStopped   = false
	ethStopped          = false
	kafkaConsumerClosed = false
)

func GetBridgeStopSignal() bool {
	return stopBridge
}

func SetBridgeStopSignal(value bool) {
	if !seal {
		stopBridge = value
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
