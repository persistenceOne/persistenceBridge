package kafka

import (
	"context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/relayer/relayer"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	db2 "github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/handler"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"log"
	"time"
)

// KafkaClose: closes all kafka connections
func KafkaClose(kafkaState utils.KafkaState, end, ended chan bool) func() {
	return func() {
		end <- true
		end <- true
		end <- true
		_ = <-ended
		_ = <-ended
		_ = <-ended
		log.Println("closing all kafka clients")
		if err := kafkaState.Producer.Close(); err != nil {
			log.Println("Error in closing producer:", err)
		}
		for _, consumerGroup := range kafkaState.ConsumerGroup {
			if err := consumerGroup.Close(); err != nil {
				log.Println("Error in closing partition:", err)
			}
		}
		if err := kafkaState.Admin.Close(); err != nil {
			log.Println("Error in closing admin:", err)
		}

	}
}

// KafkaRoutine : starts kafka in a separate goRoutine, consumers will each start in different go routines
// no need to store any db, producers and consumers are inside kafkaState struct.
// use kafka.ProducerDeliverMessage() -> to produce message
// use kafka.TopicConsumer -> to consume messages.
func KafkaRoutine(kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec, chain *relayer.Chain, ethereumClient *ethclient.Client, end, ended chan bool) {
	ctx := context.Background()

	go consumeToEthMsgs(ctx, kafkaState, protoCodec, chain, ethereumClient, end, ended)
	go consumeUnbondings(ctx, kafkaState, protoCodec, chain, ethereumClient, end, ended)
	go consumeToTendermintMessages(ctx, kafkaState, protoCodec, chain, ethereumClient, end, ended)

	log.Println("started consumers")
}

func consumeToEthMsgs(ctx context.Context, state utils.KafkaState,
	protoCodec *codec.ProtoCodec, chain *relayer.Chain, ethereumClient *ethclient.Client, end, ended chan bool) {
	consumerGroup := state.ConsumerGroup[utils.GroupToEth]
	for {
		msgHandler := handler.MsgHandler{ProtoCodec: protoCodec,
			Chain: chain, EthClient: ethereumClient, Count: 0}
		err := consumerGroup.Consume(ctx, []string{utils.ToEth}, msgHandler)
		if err != nil {
			log.Println("Error in consumer group.Consume", err)
		}
		select {
		case <-end:
			log.Println("Stopping ToEth Consumer!!!")
			ended <- true
			return
		default:
			log.Println("Next Routine Eth")

		}
	}
}

func consumeToTendermintMessages(ctx context.Context, state utils.KafkaState,
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
			log.Println("Error in consumer groupRedelegate.Consume", err)
		}
		err = groupMsgUnbond.Consume(ctx, []string{utils.MsgUnbond}, msgHandler)
		if err != nil {
			log.Println("Error in consumer groupMsgUnbond.Consume for MsgUnbond", err)
		}
		err = groupMsgDelegate.Consume(ctx, []string{utils.MsgDelegate}, msgHandler)
		if err != nil {
			log.Println("Error in consumer groupMsgDelegate.Consume", err)
		}
		err = groupMsgSend.Consume(ctx, []string{utils.MsgSend}, msgHandler)
		if err != nil {
			log.Println("Error in consumer groupMsgSend.Consume", err)
		}

		err = groupRetryTendermint.Consume(ctx, []string{utils.RetryTendermint}, msgHandler)
		if err != nil {
			log.Println("Error in consumer groupRetryTendermint.Consume", err)
		}

		err = groupMsgToTendermint.Consume(ctx, []string{utils.ToTendermint}, msgHandler)
		if err != nil {
			log.Println("Error in consumer groupMsgToTendermint.Consume", err)
		}
		select {
		case <-end:
			log.Println("Stopping To-Tendermint Consumer!!!")
			ended <- true
			return
		default:
			log.Println("Next Routine Tendermint")

		}
	}
}

func consumeUnbondings(ctx context.Context, state utils.KafkaState,
	protoCodec *codec.ProtoCodec, chain *relayer.Chain, ethereumClient *ethclient.Client, end, ended chan bool) {
	ethUnbondConsumerGroup := state.ConsumerGroup[utils.GroupEthUnbond]
	for {
		nextEpochTime, err := db2.GetUnboundEpochTime()
		if err != nil {
			log.Fatalln(err)
		}
		if time.Now().Unix() > nextEpochTime.Epoch {
			msgHandler := handler.MsgHandler{ProtoCodec: protoCodec,
				Chain: chain, EthClient: ethereumClient, Count: 0}
			err := ethUnbondConsumerGroup.Consume(ctx, []string{utils.EthUnbond}, msgHandler)
			if err != nil {
				log.Println("Error in consumer group.Consume for EthUnbond ", err)
			}

			err = db2.SetUnboundEpochTime(time.Now().Add(configuration.GetAppConfig().Kafka.EthUnbondCycleTime).Unix())
			if err != nil {
				log.Fatalln(err)
			}
			ticker := time.Tick(configuration.GetAppConfig().Kafka.EthUnbondCycleTime)
			select {
			case <-end:
				log.Println("Stopping Unbondings Consumer!!!")
				ended <- true
				return
			case <-ticker:
				log.Println("Next Routine Unbond")
			}
		} else {
			ticker := time.Tick(time.Duration(nextEpochTime.Epoch - time.Now().UnixNano()))
			select {
			case <-end:
				log.Println("Stopping Unbondings Consumer!!!")
				ended <- true
				return
			case <-ticker:
				log.Println("Next Routine Unbond")

			}
		}

	}
}
