package handler

import (
	"encoding/hex"
	"errors"
	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/outgoingTx"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
	"time"
)

func (m MsgHandler) HandleToTendermint(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var kafkaMsgs []sarama.ConsumerMessage
	claimMsgChan := claim.Messages()
	ticker := time.NewTicker(configuration.GetAppConfig().Kafka.ToTendermint.Ticker)
	defer ticker.Stop()
	var kafkaMsg *sarama.ConsumerMessage
	var ok bool
ConsumerLoop:
	for {
		select {
		case <-ticker.C:
			break ConsumerLoop
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

	// 1.add to database
	var msgBytes [][]byte
	for _, msg := range kafkaMsgs {
		msgBytes = append(msgBytes, msg.Value)
	}
	index, err := db.AddKafkaTendermintConsume(kafkaMsg.Offset, msgBytes)
	if err != nil {
		return err
	}
	// 2.set kafka offset so all next steps are independent of kafka consumer queue.
	session.MarkMessage(kafkaMsg, "")

	err = SendBatchToTendermint(index, m)
	if err != nil {
		return err
	}
	return nil
}

func ConvertMsgBytesToSDKMsg(msgBytes [][]byte, protoCodec *codec.ProtoCodec) ([]sdk.Msg, error) {
	var msgs []sdk.Msg
	for _, msgByte := range msgBytes {
		var msg sdk.Msg
		err := protoCodec.UnmarshalInterface(msgByte, &msg)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

// SendBatchToTendermint :
func SendBatchToTendermint(index uint64, handler MsgHandler) error {
	kafkaConsume, err := db.GetKafkaTendermintConsume(index)
	if err != nil {
		logging.Fatal(err)
	}
	msgs, err := ConvertMsgBytesToSDKMsg(kafkaConsume.MsgBytes, handler.ProtoCodec)
	if err != nil {
		logging.Fatal(err)
	}

	countPendingTx, err := db.CountTotalOutgoingTendermintTx()
	if err != nil {
		logging.Fatal(err)
	}

	attempts := 0
	txBroadcastSuccess := false
	for {
		attempts++
		if countPendingTx == 0 {
			response, err := outgoingTx.LogMessagesAndBroadcast(handler.Chain, msgs, 0)
			if err != nil {
				logging.Error("Unable to broadcast tendermint messages:", err)
				break
			}
			hexBytes, err := hex.DecodeString(response.TxHash)
			if err != nil {
				logging.Fatal(err)
			}
			err = db.UpdateKafkaTendermintConsumeTxHash(uint64(index), hexBytes)
			if err != nil {
				logging.Fatal(err)
			}

			logging.Info("Broadcast Tendermint Tx Hash:", response.TxHash)
			txBroadcastSuccess = true
			err = db.SetOutgoingTendermintTx(db.NewOutgoingTMTransaction(response.TxHash))
			if err != nil {
				logging.Fatal(err)
			}
			break
		} else {
			logging.Info("cannot broadcast yet, tendermint txs pending")
			time.Sleep(configuration.GetAppConfig().Tendermint.AvgBlockTime)
			countPendingTx, err = db.CountTotalOutgoingTendermintTx()
			if err != nil {
				logging.Fatal(err)
			}
		}
		if attempts >= configuration.GetAppConfig().Kafka.MaxTendermintTxAttempts {
			logging.Error("Unable to broadcast tendermint messages, max attempts while waiting for previous tx to be finished")
			break
		}
	}
	if !txBroadcastSuccess {
		addToRetryTendermintQueue(msgs, index, handler)
	}
	return nil
}

func addToRetryTendermintQueue(msgs []sdk.Msg, index uint64, handler MsgHandler) {
	config := utils.SaramaConfig()
	producer := utils.NewProducer(configuration.GetAppConfig().Kafka.Brokers, config)
	defer func() {
		err := producer.Close()
		if err != nil {
			logging.Error("failed to close producer in topic: SendBatchToTendermint, error:", err)
		}
	}()
	err := db.DeleteKafkaTendermintConsume(index)
	if err != nil {
		logging.Error("Failed to delete Tendermint msg at index: ", index, " Error: ", err)
	}
	for _, msg := range msgs {
		if msg.Type() != distributionTypes.TypeMsgWithdrawDelegatorReward {
			msgBytes, err := handler.ProtoCodec.MarshalInterface(msg)
			if err != nil {
				logging.Error("Retry txs: Failed to Marshal ToTendermint Retry msg:", msg.String(), "Error:", err)
			}
			err = utils.ProducerDeliverMessage(msgBytes, utils.RetryTendermint, producer)
			if err != nil {
				logging.Error("Retry txs: Failed to add msg to kafka RetryTendermint queue, Msg:", msg.String(), "Error:", err)
			}
			logging.Info("Retry txs: Produced to kafka for topic RetryTendermint:", msg.String())
		}
	}

}
