package constants

import (
	"github.com/Shopify/sarama"
	"time"
)

var (
	DefaultBrokers            = []string{"localhost:9092"}
	EthBatchSize              = 2
	TendermintBatchSize       = 3
	DefaultEthUnbondCycleTime = time.Duration(259200000000000) //3days

	// TopicDetail : configs only required for admin to create topics if not present.
	TopicDetail = sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
)
