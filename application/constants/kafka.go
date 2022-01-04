/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package constants

import (
	"time"
)

var (
	DefaultBrokers            = "localhost:9092"
	MinEthBatchSize           = 1
	MaxEthBatchSize           = 30
	EthTicker                 = 30 * time.Second
	MinTendermintBatchSize    = 1 //Do not change
	MaxTendermintBatchSize    = 30
	TendermintTicker          = 3 * time.Second
	DefaultEthUnbondCycleTime = 259200 * time.Second //3days in seconds

	// TopicDetail: configs only required for admin to create topics if not present.
	TopicDetailNumPartitions     = 1
	TopicDetailReplicationFactor = 1
)
