/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package utils

import "github.com/Shopify/sarama"

func SaramaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_1_0                 // hardcoded
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 3                    // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true

	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	return config
}
