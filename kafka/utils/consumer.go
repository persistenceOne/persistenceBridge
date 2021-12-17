/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"github.com/Shopify/sarama"
)

func NewConsumerGroup(kafkaPorts []string, groupID string, config *sarama.Config) sarama.ConsumerGroup {
	consumerGroup, err := sarama.NewConsumerGroup(kafkaPorts, groupID, config)
	if err != nil {
		panic(err)
	}

	return consumerGroup
}
