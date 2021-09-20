package handler

import (
	"encoding/json"
	"errors"
	"github.com/Shopify/sarama"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/outgoingTx"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
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
			} else {
				return nil
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
	msgs := make([]outgoingTx.WrapTokenMsg, len(kafkaMsgs))
	for i, kafkaMsg := range kafkaMsgs {
		var msg outgoingTx.WrapTokenMsg
		err := json.Unmarshal(kafkaMsg.Value, &msg)
		if err != nil {
			return nil, err
		}
		msgs[i] = msg
	}
	return msgs, nil
}

// SendBatchToEth : Handling of msgSend
func SendBatchToEth(kafkaMsgs []sarama.ConsumerMessage, handler MsgHandler) error {
	msgs, err := convertKafkaMsgsToEthMsg(kafkaMsgs)
	if err != nil {
		return err
	}
	logging.Info("batched messages to send to ETH:", msgs)

	hash, err := outgoingTx.EthereumWrapToken(handler.EthClient, msgs)
	if err != nil {
		logging.Error("Unable to do ethereum tx (adding messages again to kafka), messages:", msgs, "error:", err)
		config := utils.SaramaConfig()
		producer := utils.NewProducer(configuration.GetAppConfig().Kafka.Brokers, config)
		defer func() {
			err := producer.Close()
			if err != nil {
				logging.Error("failed to close producer in topic: SendBatchToEth, err:", err)
			}
		}()

		for i, kafkaMsg := range kafkaMsgs {
			err = utils.ProducerDeliverMessage(kafkaMsg.Value, utils.ToEth, producer)
			if err != nil {
				logging.Error("Failed to add msg to kafka queue, message:", msgs[i], "error:", err)
				// TODO @Puneet continue or return? ~ Log (ALERT) and continue, need to manually do the failed ones.
			}
		}
		return err
	} else {
		err = db.SetOutgoingEthereumTx(db.NewOutgoingETHTransaction(hash, msgs))
		if err != nil {
			logging.Fatal(err)
		}
	}
	logging.Info("Broadcast Eth Tx hash:", hash.String())
	return nil
}
