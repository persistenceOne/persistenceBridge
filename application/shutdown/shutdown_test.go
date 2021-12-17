/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package shutdown

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetBridgeStopSignal(t *testing.T) {
	require.Equal(t, stopBridge, GetBridgeStopSignal())
}

func TestSetBridgeStopSignal(t *testing.T) {
	SetBridgeStopSignal(true)
	require.Equal(t, true, seal)
	require.Equal(t, true, GetBridgeStopSignal())
}

func TestGetTMStopped(t *testing.T) {
	require.Equal(t, tendermintStopped, GetTMStopped())
}

func TestSetTMStopped(t *testing.T) {
	SetTMStopped(true)
	require.Equal(t, true, GetTMStopped())
}

func TestGetEthStopped(t *testing.T) {
	require.Equal(t, ethStopped, GetETHStopped())
}

func TestSetEthStopped(t *testing.T) {
	SetETHStopped(true)
	require.Equal(t, true, GetETHStopped())
}

func TestGetKafkaConsumerClosed(t *testing.T) {
	require.Equal(t, kafkaConsumerClosed, GetKafkaConsumerClosed())
}

func TestSetKafkaConsumerClosed(t *testing.T) {
	SetKafkaConsumerClosed(true)
	require.Equal(t, true, GetKafkaConsumerClosed())
}
