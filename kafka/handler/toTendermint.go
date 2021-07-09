package handler

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/outgoingTx"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"log"
	"time"
)

func (m MsgHandler) HandleToTendermint(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var kafkaMsgs []sarama.ConsumerMessage
	claimMsgChan := claim.Messages()
	ticker := time.Tick(configuration.GetAppConfig().Kafka.ToTendermint.Ticker)
	var kafkaMsg *sarama.ConsumerMessage
	var ok bool
ConsumerLoop:
	for {
		select {
		case <-ticker:
			if len(kafkaMsgs) >= configuration.GetAppConfig().Kafka.ToTendermint.MinBatchSize {
				break ConsumerLoop
			}
		case kafkaMsg, ok = <-claimMsgChan:
			if ok {
				kafkaMsgs = append(kafkaMsgs, *kafkaMsg)
				if len(kafkaMsgs) == configuration.GetAppConfig().Kafka.ToTendermint.MaxBatchSize {
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
	err := SendBatchToTendermint(kafkaMsgs, m)
	if err != nil {
		return err
	}
	session.MarkMessage(kafkaMsg, "")
	return nil
}

func ConvertKafkaMsgsToSDKMsg(kafkaMsgs []sarama.ConsumerMessage, protoCodec *codec.ProtoCodec) ([]sdk.Msg, error) {
	var msgs []sdk.Msg
	for _, kafkaMsg := range kafkaMsgs {
		var msg sdk.Msg
		err := protoCodec.UnmarshalInterface(kafkaMsg.Value, &msg)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

// SendBatchToTendermint :
func SendBatchToTendermint(kafkaMsgs []sarama.ConsumerMessage, handler MsgHandler) error {
	msgs, err := ConvertKafkaMsgsToSDKMsg(kafkaMsgs, handler.ProtoCodec)
	if err != nil {
		return err
	}

	// TODO set memo and timeout height
	response, err := outgoingTx.FilterMessagesAndBroadcast(handler.Chain, msgs, 0)
	if err != nil {
		log.Printf("error occured while send to Tendermint:%v\n", err)
		config := utils.SaramaConfig()
		producer := utils.NewProducer(configuration.GetAppConfig().Kafka.Brokers, config)
		defer func() {
			err := producer.Close()
			if err != nil {
				log.Printf("failed to close producer in topic: SendBatchToTendermint\n")
			}
		}()

		for _, msg := range msgs {
			if msg.Type() != distributionTypes.TypeMsgWithdrawDelegatorReward {
				msgBytes, err := handler.ProtoCodec.MarshalInterface(msg)
				if err != nil {
					log.Printf("Retry txs: Failed to Marshal ToTendermint Retry msg: Error: %v\n", err)
				}
				err = utils.ProducerDeliverMessage(msgBytes, utils.ToTendermint, producer)
				if err != nil {
					log.Printf("Retry txs: Failed to add msg to kafka queue: %s\n", err.Error())
				}
				log.Printf("Retry txs: Produced to kafka: %v, for topic %v\n", msg.Type(), utils.ToTendermint)
			}
		}
		return nil
	} else {
		err = db.SetTendermintTx(db.NewTMTransaction(response.TxHash))
		if err != nil {
			panic(err)
		}
	}
	log.Printf("Broadcasted Tendermint TX HASH: %s\n", response.TxHash)
	return nil
}
