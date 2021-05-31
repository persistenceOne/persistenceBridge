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

	// These are the configuration parameters for running kafka admins and producers and consumers. Declared very minimal
	replicaAssignment = map[int32][]int32{}
	configEntries     = map[string]*string{}

	// TopicDetail : configs only required for admin to create topics if not present.
	TopicDetail = sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
		ReplicaAssignment: replicaAssignment,
		ConfigEntries:     configEntries,
	}
)
