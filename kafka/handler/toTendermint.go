/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package handler

import (
	"time"

	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/outgoingtx"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func (m MsgHandler) HandleToTendermint(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var kafkaMsgs []sarama.ConsumerMessage

	claimMsgChan := claim.Messages()

	ticker := time.NewTicker(configuration.GetAppConfig().Kafka.ToTendermint.Ticker)
	defer ticker.Stop()

	var (
		kafkaMsg *sarama.ConsumerMessage
		ok       bool
	)

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
		return ErrKafkaNilMessage
	}

	err := SendBatchToTendermint(kafkaMsgs, m)
	if err != nil {
		return err
	}

	session.MarkMessage(kafkaMsg, "")

	return nil
}

func ConvertKafkaMsgsToSDKMsg(kafkaMsgs []sarama.ConsumerMessage, protoCodec *codec.ProtoCodec) ([]sdk.Msg, error) {
	msgs := make([]sdk.Msg, len(kafkaMsgs))

	for i := range kafkaMsgs {
		var msg sdk.Msg

		err := protoCodec.UnmarshalInterface(kafkaMsgs[i].Value, &msg)
		if err != nil {
			return nil, err
		}

		msgs[i] = msg
	}

	return msgs, nil
}

// SendBatchToTendermint :
func SendBatchToTendermint(kafkaMsgs []sarama.ConsumerMessage, handler MsgHandler) error {
	msgs, err := ConvertKafkaMsgsToSDKMsg(kafkaMsgs, handler.ProtoCodec)
	if err != nil {
		return err
	}

	countPendingTx, err := db.CountTotalOutgoingTendermintTx()
	if err != nil {
		logging.Fatal(err)
	}

	for {
		if countPendingTx == 0 {
			var response *sdk.TxResponse

			response, err = outgoingtx.LogMessagesAndBroadcast(handler.Chain, msgs, 0)
			if err != nil {
				logging.Error("Unable to broadcast tendermint messages:", err)

				config := utils.SaramaConfig()
				producer := utils.NewProducer(configuration.GetAppConfig().Kafka.Brokers, config)

				func() {
					defer func() {
						innerErr := producer.Close()
						if innerErr != nil {
							logging.Error("failed to close producer in topic: SendBatchToTendermint, error:", innerErr)
						}
					}()

					for _, msg := range msgs {
						if msg.Type() == distributionTypes.TypeMsgWithdrawDelegatorReward {
							continue
						}

						var msgBytes []byte

						msgBytes, err = handler.ProtoCodec.MarshalInterface(msg)
						if err != nil {
							logging.Error("Retry txs: Failed to Marshal ToTendermint Retry msg:", msg.String(), "Error:", err)
							// TODO @Puneet continue or return? ~ this case should never come, log(ALERT), continue
						}

						err = utils.ProducerDeliverMessage(msgBytes, utils.RetryTendermint, producer)
						if err != nil {
							logging.Error("Retry txs: Failed to add msg to kafka queue, Msg:", msg.String(), "Error:", err)
							// TODO @Puneet continue or return? ~ let it continue, log the message, will have to send manually.
						}

						logging.Info("Retry txs: Produced to kafka for topic RetryTendermint:", msg.String())
					}
				}()

				return nil
			}

			err = db.SetOutgoingTendermintTx(db.NewOutgoingTMTransaction(response.TxHash))
			if err != nil {
				logging.Fatal(err)
			}

			logging.Info("Broadcast Tendermint Tx Hash:", response.TxHash)

			return nil
		}

		logging.Info("cannot broadcast yet, tendermint txs pending")

		time.Sleep(4 * time.Second)

		countPendingTx, err = db.CountTotalOutgoingTendermintTx()
		if err != nil {
			logging.Fatal(err)
		}
	}
}
