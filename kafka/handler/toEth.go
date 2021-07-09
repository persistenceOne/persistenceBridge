package handler

import (
	"encoding/json"
	"errors"
	"github.com/Shopify/sarama"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/outgoingTx"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"log"
	"time"
)

func (m MsgHandler) HandleToEth(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var kafkaMsgs []sarama.ConsumerMessage
	claimMsgChan := claim.Messages()
	ticker := time.Tick(configuration.GetAppConfig().Kafka.ToEth.Ticker)
	var kafkaMsg *sarama.ConsumerMessage
	var ok bool
ConsumerLoop:
	for {
		select {
		case <-ticker:
			if len(kafkaMsgs) >= configuration.GetAppConfig().Kafka.ToEth.MinBatchSize {
				break ConsumerLoop
			}
		case kafkaMsg, ok = <-claimMsgChan:
			if ok {
				kafkaMsgs = append(kafkaMsgs, *kafkaMsg)
				if len(kafkaMsgs) == configuration.GetAppConfig().Kafka.ToEth.MaxBatchSize {
					break ConsumerLoop
				}
			} else {
				break ConsumerLoop
			}
		}
	}

	if len(kafkaMsgs) == 0 {
		return nil
	}
	if kafkaMsg == nil {
		return errors.New("kafka returned nil message")
	}
	err := SendBatchToEth(kafkaMsgs, m)
	if err != nil {
		return err
	}
	session.MarkMessage(kafkaMsg, "")
	return nil
}

func convertKafkaMsgsToEthMsg(kafkaMsgs []sarama.ConsumerMessage) ([]outgoingTx.WrapTokenMsg, error) {
	var msgs []outgoingTx.WrapTokenMsg
	for _, kafkaMsg := range kafkaMsgs {
		var msg outgoingTx.WrapTokenMsg
		err := json.Unmarshal(kafkaMsg.Value, &msg)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

// SendBatchToEth : Handling of msgSend
func SendBatchToEth(kafkaMsgs []sarama.ConsumerMessage, handler MsgHandler) error {
	msgs, err := convertKafkaMsgsToEthMsg(kafkaMsgs)
	if err != nil {
		return err
	}
	log.Printf("batched messages to send to ETH: %v\n", msgs)

	hash, err := outgoingTx.EthereumWrapToken(handler.EthClient, msgs)
	if err != nil {
		log.Printf("error occuerd in sending eth transaction: %v, adding messages agin to kafka\n", err)
		config := utils.SaramaConfig()
		producer := utils.NewProducer(configuration.GetAppConfig().Kafka.Brokers, config)
		defer func() {
			err := producer.Close()
			if err != nil {
				log.Printf("failed to close producer in topic: SendBatchToEth\n")
			}
		}()

		for _, kafkaMsg := range kafkaMsgs {
			err = utils.ProducerDeliverMessage(kafkaMsg.Value, utils.ToEth, producer)
			if err != nil {
				log.Printf("Failed to add msg to kafka queue: %s\n", err.Error())
			}
		}
		return err
	} else {
		err = db.SetEthereumTx(db.NewETHTransaction(hash, msgs))
		if err != nil {
			panic(err)
		}
	}
	log.Printf("Broadcasted Eth Tx hash: %s\n", hash.String())
	return nil
}
