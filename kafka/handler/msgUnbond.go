package handler

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"log"
)

func (m MsgHandler) HandleMsgUnbond(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	config := utils.SaramaConfig()
	producer := utils.NewProducer(configuration.GetAppConfig().Kafka.Brokers, config)
	defer func() {
		err := producer.Close()
		if err != nil {
			log.Printf("failed to close producer in topic: %v\n", utils.MsgUnbond)
		}
	}()

	claimMsgChan := claim.Messages()
	var kafkaMsg *sarama.ConsumerMessage
	var ok bool
ConsumerLoop:
	for {
		select {
		case kafkaMsg, ok = <-claimMsgChan:
			if !ok {
				break ConsumerLoop
			}
			if kafkaMsg == nil {
				return errors.New("kafka returned nil message")
			}
			err := utils.ProducerDeliverMessage(kafkaMsg.Value, utils.ToTendermint, producer)
			if err != nil {
				log.Printf("failed to produce from %v to :%v", utils.MsgUnbond, utils.ToTendermint)
				break ConsumerLoop
			}
			session.MarkMessage(kafkaMsg, "")
			m.Count++
		default:
			break ConsumerLoop
		}
	}

	return nil
}
