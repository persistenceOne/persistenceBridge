package handler

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"log"
)

func (m MsgHandler) HandleMsgDelegate(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	config := utils.SaramaConfig()
	producer := utils.NewProducer(m.PstakeConfig.Kafka.Brokers, config)
	defer func() {
		err := producer.Close()
		if err != nil {
			log.Printf("failed to close producer in topic: %v\n", utils.MsgDelegate)
		}
	}()
	messagesLength := len(claim.Messages())
	if messagesLength > 0 {
		var msgs [][]byte
		var kafkaMsg *sarama.ConsumerMessage
		for i := 0; i < messagesLength; i++ {
			kafkaMsg = <-claim.Messages()
			if kafkaMsg == nil {
				return errors.New("kafka returned nil message")
			}

			msgs = append(msgs, kafkaMsg.Value)
		}
		if len(msgs) > 0 {
			//TODO send as one msgDelegate
			err := utils.ProducerDeliverMessages(msgs, utils.ToTendermint, producer)
			session.MarkMessage(kafkaMsg, "")
			if err != nil {
				log.Printf("error in handler for topic %v, failed to produce to queue\n", utils.MsgDelegate)
				return err
			}
		}
	}
	m.Count += messagesLength
	return nil
}
