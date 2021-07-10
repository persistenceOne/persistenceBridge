package constants

import (
	"github.com/Shopify/sarama"
	"time"
)

var (
	DefaultBrokers            = []string{"localhost:9092"}
	MinEthBatchSize           = 1
	MaxEthBatchSize           = 4
	EthTicker                 = 5 * time.Second
	MinTendermintBatchSize    = 1
	MaxTendermintBatchSize    = 5
	TendermintTicker          = 5 * time.Second
	DefaultEthUnbondCycleTime = 259200 * time.Second //3days in seconds

	// TopicDetail : configs only required for admin to create topics if not present.
	TopicDetail = sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
)
