package constants

import (
	"github.com/Shopify/sarama"
	"time"
)

var (
	DefaultBrokers            = []string{"localhost:9092"}
	MinEthBatchSize           = 1
	MaxEthBatchSize           = 30
	EthTicker                 = 30 * time.Second
	MinTendermintBatchSize    = 1 //Do not change
	MaxTendermintBatchSize    = 30
	TendermintTicker          = 3 * time.Second
	DefaultEthUnbondCycleTime = 259200 * time.Second //3days in seconds

	// TopicDetail : configs only required for admin to create topics if not present.
	TopicDetail = sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
)
