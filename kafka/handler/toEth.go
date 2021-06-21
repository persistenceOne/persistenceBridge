package handler

import (
	"encoding/json"
	"errors"
	"github.com/Shopify/sarama"
	"github.com/persistenceOne/persistenceBridge/application"
	ethereum2 "github.com/persistenceOne/persistenceBridge/ethereum"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"log"
	"time"
)

func (m MsgHandler) HandleToEth(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim, batchSize int) error {
	var kafkaMsgs []sarama.ConsumerMessage
	claimMsgChan := claim.Messages()
	closeChan := make(chan bool, 1)
	ticker := time.Tick(1 * time.Second)
	var kafkaMsg *sarama.ConsumerMessage
	var ok bool
	for {
		select {
		case <-closeChan:
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
		case <-ticker:
			if len(kafkaMsgs) != 0 {
				AddToBufferedChannelIfCapacityPermits(closeChan, true)
			}
		case kafkaMsg, ok = <-claimMsgChan:
			if ok {
				kafkaMsgs = append(kafkaMsgs, *kafkaMsg)
				if len(kafkaMsgs) == batchSize {
					AddToBufferedChannelIfCapacityPermits(closeChan, true)
				} else if len(kafkaMsgs) > batchSize {
					log.Printf("Select tried to batch more messages in handler: %v ,not "+
						"comitting offset, %v", utils.ToEth, kafkaMsg.Offset)
					return nil
				}
			} else {
				AddToBufferedChannelIfCapacityPermits(closeChan, true)
			}
		}
	}

	//msgs := make([]sarama.ConsumerMessage, 0, batchSize)
	//for {
	//	kafkaMsg := <-claim.Messages()
	//	if kafkaMsg == nil {
	//		return errors.New("kafka returned nil message")
	//	}
	//	log.Printf("Message topic:%q partition:%d offset:%d\n", kafkaMsg.Topic, kafkaMsg.Partition, kafkaMsg.Offset)
	//
	//	ok, err := BatchAndHandleEthereum(&msgs, *kafkaMsg, m)
	//	if ok && err == nil {
	//		session.MarkMessage(kafkaMsg, "")
	//		return nil
	//	}
	//	if err != nil {
	//		return err
	//	}
	//}
}

func ConvertKafkaMsgsToEthMsg(kafkaMsgs []sarama.ConsumerMessage) ([]ethereum2.EthTxMsg, error) {
	var msgs []ethereum2.EthTxMsg
	for _, kafkaMsg := range kafkaMsgs {
		var msg ethereum2.EthTxMsg
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
	msgs, err := ConvertKafkaMsgsToEthMsg(kafkaMsgs)
	if err != nil {
		return err
	}
	log.Printf("batched messages to send to ETH: %v\n", msgs)

	hash, err := ethereum2.SendTxToEth(handler.EthClient, msgs, application.GetAppConfiguration().EthGasLimit)
	if err != nil {
		log.Printf("error occuerd in sending eth transaction: %v\n", err)
		return err
	}
	log.Printf("Broadcasted Eth Tx hash: %s\n", hash)
	return nil
}
