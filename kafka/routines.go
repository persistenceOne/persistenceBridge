/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package kafka

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/relayer/relayer"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	db2 "github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/handler"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

// Close closes all kafka connections
func Close(kafkaState *utils.KafkaState, end, ended chan bool) func() {
	return func() {
		end <- true
		end <- true
		end <- true

		<-ended
		<-ended
		<-ended

		logging.Info("closing all kafka clients")

		if err := kafkaState.Producer.Close(); err != nil {
			logging.Error("Closing producer error:", err)
		}

		for _, consumerGroup := range kafkaState.ConsumerGroup {
			if err := consumerGroup.Close(); err != nil {
				logging.Error("Closing partition error:", err)
			}
		}

		if err := kafkaState.Admin.Close(); err != nil {
			logging.Error("Closing admin error:", err)
		}
	}
}

// Routine : starts kafka in a separate goRoutine, consumers will each start in different go routines
// no need to store any db, producers and consumers are inside kafkaState struct.
// use kafka.ProducerDeliverMessage() -> to produce message
// use kafka.TopicConsumer -> to consume messages.
func Routine(kafkaState *utils.KafkaState, protoCodec *codec.ProtoCodec, chain *relayer.Chain, ethereumClient *ethclient.Client, end, ended chan bool) {
	ctx := context.Background()

	go consumeToEthMsgs(ctx, kafkaState, protoCodec, chain, ethereumClient, end, ended)
	go consumeUnbondings(ctx, kafkaState, protoCodec, chain, ethereumClient, end, ended)
	go consumeToTendermintMessages(ctx, kafkaState, protoCodec, chain, ethereumClient, end, ended)

	logging.Info("Started consumers")
}

func consumeToEthMsgs(ctx context.Context, state *utils.KafkaState,
	protoCodec *codec.ProtoCodec, chain *relayer.Chain, ethereumClient *ethclient.Client, end, ended chan bool) {
	consumerGroup := state.ConsumerGroup[utils.GroupToEth]

	for {
		msgHandler := handler.MsgHandler{ProtoCodec: protoCodec,
			Chain: chain, EthClient: ethereumClient, Count: 0}

		err := consumerGroup.Consume(ctx, []string{utils.ToEth}, msgHandler)
		if err != nil {
			logging.Error("Consumer group.Consume:", err)
		}

		select {
		case <-end:
			logging.Info("Stopping ToEth Consumer!!!")
			ended <- true

			return
		default:
			logging.Debug("Next Routine Eth")
		}
	}
}

func consumeToTendermintMessages(ctx context.Context, state *utils.KafkaState,
	protoCodec *codec.ProtoCodec, chain *relayer.Chain, ethereumClient *ethclient.Client, end, ended chan bool) {
	groupMsgUnbond := state.ConsumerGroup[utils.GroupMsgUnbond]
	groupMsgDelegate := state.ConsumerGroup[utils.GroupMsgDelegate]
	groupMsgSend := state.ConsumerGroup[utils.GroupMsgSend]
	groupRedelegate := state.ConsumerGroup[utils.GroupRedelegate]
	groupRetryTendermint := state.ConsumerGroup[utils.GroupRetryTendermint]
	groupMsgToTendermint := state.ConsumerGroup[utils.GroupToTendermint]

	for {
		msgHandler := handler.MsgHandler{ProtoCodec: protoCodec,
			Chain: chain, EthClient: ethereumClient, Count: 0, WithdrawRewards: false}

		err := groupRedelegate.Consume(ctx, []string{utils.Redelegate}, msgHandler)
		if err != nil {
			logging.Error("Consumer groupRedelegate.Consume:", err)
		}

		err = groupMsgUnbond.Consume(ctx, []string{utils.MsgUnbond}, msgHandler)
		if err != nil {
			logging.Error("Consumer groupMsgUnbond.Consume:", err)
		}

		err = groupMsgDelegate.Consume(ctx, []string{utils.MsgDelegate}, msgHandler)
		if err != nil {
			logging.Error("Consumer groupMsgDelegate.Consume:", err)
		}

		err = groupMsgSend.Consume(ctx, []string{utils.MsgSend}, msgHandler)
		if err != nil {
			logging.Error("Consumer groupMsgSend.Consume:", err)
		}

		err = groupRetryTendermint.Consume(ctx, []string{utils.RetryTendermint}, msgHandler)
		if err != nil {
			logging.Error("Consumer groupRetryTendermint.Consume:", err)
		}

		err = groupMsgToTendermint.Consume(ctx, []string{utils.ToTendermint}, msgHandler)
		if err != nil {
			logging.Error("Consumer groupMsgToTendermint.Consume:", err)
		}

		select {
		case <-end:
			logging.Info("Stopping To-Tendermint Consumer!!!")
			ended <- true

			return
		default:
			logging.Debug("Next Routine Tendermint")
		}
	}
}

func consumeUnbondings(ctx context.Context, state *utils.KafkaState,
	protoCodec *codec.ProtoCodec, chain *relayer.Chain, ethereumClient *ethclient.Client, end, ended chan bool) {
	ethUnbondConsumerGroup := state.ConsumerGroup[utils.GroupEthUnbond]

	for {
		nextEpochTime, err := db2.GetUnboundEpochTime()
		if err != nil {
			logging.Fatal(err)
		}

		if time.Now().Unix() > nextEpochTime.Epoch {
			msgHandler := handler.MsgHandler{ProtoCodec: protoCodec,
				Chain: chain, EthClient: ethereumClient, Count: 0}

			err := ethUnbondConsumerGroup.Consume(ctx, []string{utils.EthUnbond}, msgHandler)
			if err != nil {
				logging.Error("Consumer group.Consume for EthUnbond:", err)
			}

			err = db2.SetUnboundEpochTime(nextEpochTime.Epoch + configuration.GetAppConfig().Kafka.EthUnbondCycleTime.Milliseconds()/1000)
			if err != nil {
				logging.Fatal(err)
			}

			ticker := time.NewTicker(10 * time.Second)

			select {
			case <-end:
				logging.Info("Stopping Unbondings Consumer!!!")

				ended <- true

				ticker.Stop()

				return
			case <-ticker.C:
				logging.Debug("Next Routine Unbond")

				ticker.Stop()
			}
		} else {
			ticker := time.NewTicker(10 * time.Second)

			select {
			case <-end:
				logging.Info("Stopping Unbondings Consumer!!!")

				ended <- true
				ticker.Stop()

				return
			case <-ticker.C:
				logging.Debug("Next Routine Unbond")
				ticker.Stop()
			}
		}
	}
}
