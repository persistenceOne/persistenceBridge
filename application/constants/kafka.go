/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package constants

import (
	"time"

	"github.com/Shopify/sarama"
)

const (
	DefaultBroker             = "localhost:9092"
	MinEthBatchSize           = 1
	MaxEthBatchSize           = 30
	EthTicker                 = 30 * time.Second
	MinTendermintBatchSize    = 1 // Do not change
	MaxTendermintBatchSize    = 30
	TendermintTicker          = 3 * time.Second
	DefaultEthUnbondCycleTime = 259200 * time.Second // 3days in seconds
)

var (
	// TopicDetail : configs only required for admin to create topics if not present.
	// nolint value struct is safe to be used as global variable
	// nolint: gochecknoglobals
	TopicDetail = TopicDetails{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
)

type TopicDetails struct {
	NumPartitions     int32
	ReplicationFactor int16
}

func ToKafkaTopicDetail(d TopicDetails) sarama.TopicDetail {
	return sarama.TopicDetail{
		NumPartitions:     d.NumPartitions,
		ReplicationFactor: d.ReplicationFactor,
	}
}
