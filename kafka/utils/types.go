/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"github.com/Shopify/sarama"
)

// Ticket : is a type that implements string
type Ticket string

// // KafkaMsg : is a store that can be stored in kafka queues
// type KafkaMsg struct {
//	MsgBytes      sdk.MsgBytes `json:"msg"`
//	TicketID Ticket  `json:"ticketID"`
// }
//
// // NewKafkaMsgFromRest : makes a msg to send to kafka queue
// func NewKafkaMsgFromRest(msg sdk.MsgBytes, ticketID Ticket) KafkaMsg {
//	return KafkaMsg{
//		MsgBytes:      msg,
//		TicketID: ticketID,
//	}
// }

// TicketIDResponse : is a json structure to send TicketID to user
type TicketIDResponse struct {
	TicketID Ticket `json:"ticketID" valid:"required~ticketID is mandatory,length(20)~ticketID length should be 20" `
}

// KafkaState : is a struct showing the state of kafka
type KafkaState struct {
	HomeDir       string
	Admin         sarama.ClusterAdmin
	ConsumerGroup map[string]sarama.ConsumerGroup
	Producer      sarama.SyncProducer
	Topics        []string
}

// NewKafkaState : returns a kafka state
func NewKafkaState(kafkaPorts []string, homeDir string, topicDetail sarama.TopicDetail) *KafkaState {
	config := SaramaConfig()
	admin := KafkaAdmin(kafkaPorts, config)

	adminTopics, err := admin.ListTopics()
	if err != nil {
		panic(err)
	}

	// create topics if not present
	for _, topic := range Topics() {
		if _, ok := adminTopics[topic]; !ok {
			TopicsInit(admin, topic, topicDetail)
		}
	}

	producer := NewProducer(kafkaPorts, config)
	groups := Groups()

	consumers := make(map[string]sarama.ConsumerGroup, len(groups))
	for _, group := range groups {
		consumers[group] = NewConsumerGroup(kafkaPorts, group, config)
	}

	return &KafkaState{
		HomeDir:       homeDir,
		Admin:         admin,
		ConsumerGroup: consumers,
		Producer:      producer,
		Topics:        Topics(),
	}
}
