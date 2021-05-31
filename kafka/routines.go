package kafka

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/relayer/relayer"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/kafka/handler"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"log"
	"time"
)

// KafkaClose: closes all kafka connections
func KafkaClose(kafkaState utils.KafkaState) func() {
	return func() {
		fmt.Println("closing all kafka clients.")
		if err := kafkaState.Producer.Close(); err != nil {
			log.Print("Error in closing producer:", err)
		}
		for _, consumerGroup := range kafkaState.ConsumerGroup {
			if err := consumerGroup.Close(); err != nil {
				log.Print("Error in closing partition:", err)
			}
		}
		if err := kafkaState.Admin.Close(); err != nil {
			log.Print("Error in closing admin:", err)
		}

	}
}

// KafkaRoutine: starts kafka in a separate goRoutine, consumers will each start in different go routines
// no need to store any db, producers and consumers are inside kafkaState struct.
// use kafka.ProducerDeliverMessage() -> to produce message
// use kafka.TopicConsumer -> to consume messages.
func KafkaRoutine(kafkaState utils.KafkaState, pstakeConfig configuration.Config, protoCodec *codec.ProtoCodec, chain *relayer.Chain, ethereumClient *ethclient.Client) {
	ctx := context.Background()

	go consumeToEthMsgs(ctx, kafkaState, pstakeConfig, protoCodec, chain, ethereumClient)
	go consumeUnbondings(ctx, kafkaState, pstakeConfig, protoCodec, chain, ethereumClient)
	go consumeToTendermintMessages(ctx, kafkaState, pstakeConfig, protoCodec, chain, ethereumClient)
	// go consume other messages

	fmt.Println("started consumers")
}

func consumeToEthMsgs(ctx context.Context, state utils.KafkaState, kafkaConfig configuration.Config,
	protoCodec *codec.ProtoCodec, chain *relayer.Chain, ethereumClient *ethclient.Client) {
	consumerGroup := state.ConsumerGroup[utils.GroupToEth]
	for {
		msgHandler := handler.MsgHandler{PstakeConfig: kafkaConfig, ProtoCodec: protoCodec,
			Chain: chain, EthClient: ethereumClient, Count: 0}
		err := consumerGroup.Consume(ctx, []string{utils.ToEth}, msgHandler)
		if err != nil {
			log.Println("Error in consumer group.Consume", err)
		}
		time.Sleep(time.Duration(1000000000))
	}
}

func consumeToTendermintMessages(ctx context.Context, state utils.KafkaState, pstakeConfig configuration.Config,
	protoCodec *codec.ProtoCodec, chain *relayer.Chain, ethereumClient *ethclient.Client) {
	groupMsgUnbond := state.ConsumerGroup[utils.GroupMsgUnbond]
	groupMsgDelegate := state.ConsumerGroup[utils.GroupMsgDelegate]
	groupMsgSend := state.ConsumerGroup[utils.GroupMsgSend]
	groupMsgToTendermint := state.ConsumerGroup[utils.GroupToTendermint]
	for {
		msgHandler := handler.MsgHandler{PstakeConfig: pstakeConfig, ProtoCodec: protoCodec,
			Chain: chain, EthClient: ethereumClient, Count: 0}
		err := groupMsgUnbond.Consume(ctx, []string{utils.MsgUnbond}, msgHandler)
		if err != nil {
			log.Println("Error in consumer group.Consume for MsgUnbond", err)
		}
		err = groupMsgDelegate.Consume(ctx, []string{utils.MsgDelegate}, msgHandler)
		if err != nil {
			log.Println("Error in consumer group.Consume", err)
		}
		err = groupMsgSend.Consume(ctx, []string{utils.MsgSend}, msgHandler)
		if err != nil {
			log.Println("Error in consumer group.Consume", err)
		}
		err = groupMsgToTendermint.Consume(ctx, []string{utils.ToTendermint}, msgHandler)
		if err != nil {
			log.Println("Error in consumer group.Consume", err)
		}
		time.Sleep(time.Duration(1000000000))
	}
}

func consumeUnbondings(ctx context.Context, state utils.KafkaState, pstakeConfig configuration.Config,
	protoCodec *codec.ProtoCodec, chain *relayer.Chain, ethereumClient *ethclient.Client) {
	ethUnbondConsumerGroup := state.ConsumerGroup[utils.GroupEthUnbond]
	for {
		msgHandler := handler.MsgHandler{PstakeConfig: pstakeConfig, ProtoCodec: protoCodec,
			Chain: chain, EthClient: ethereumClient, Count: 0}
		err := ethUnbondConsumerGroup.Consume(ctx, []string{utils.EthUnbond}, msgHandler)
		if err != nil {
			log.Println("Error in consumer group.Consume for EthUnbond ", err)
		}
		time.Sleep(pstakeConfig.Kafka.EthUnbondCycleTime)
	}
}
