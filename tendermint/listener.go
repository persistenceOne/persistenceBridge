package tendermint

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/relayer/relayer"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/shutdown"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func StartListening(initClientCtx client.Context, chain *relayer.Chain, brokers []string, protoCodec *codec.ProtoCodec, sleepDuration time.Duration) {
	ctx := context.Background()
	kafkaProducer := utils.NewProducer(brokers, utils.SaramaConfig())
	defer func(kafkaProducer sarama.SyncProducer) {
		err := kafkaProducer.Close()
		if err != nil {
			logging.Error(err)
		}
	}(kafkaProducer)

	for {
		// For Tendermint, we can directly query without waiting for blocks since there is finality
		err := onNewBlock(ctx, initClientCtx, chain, &kafkaProducer, protoCodec)
		if err != nil {
			logging.Error("Stopping Tendermint Listener, onNewBlock err:", err)
			return
		}

		if shutdown.GetBridgeStopSignal() {
			if shutdown.GetKafkaConsumerClosed() {
				logging.Info("Stopping Tendermint Listener!!!")
				shutdown.SetTMStopped(true)
				return
			}
			time.Sleep(5 * time.Second)
			continue
		}

		abciInfo, err := chain.Client.ABCIInfo(ctx)
		if err != nil {
			logging.Error("Unable to fetch tendermint ABCI info:", err)
			time.Sleep(sleepDuration)
			continue
		}

		cosmosStatus, err := db.GetCosmosStatus()
		if err != nil {
			panic(err)
		}

		if abciInfo.Response.LastBlockHeight > cosmosStatus.LastCheckHeight {
			processHeight := cosmosStatus.LastCheckHeight + 1
			logging.Info("Tendermint Block:", processHeight)

			//TODO bug of pages and perPage
			txSearchResult, err := chain.Client.TxSearch(ctx, fmt.Sprintf("tx.height=%d", processHeight), true, nil, nil, "asc")
			if err != nil {
				logging.Error("Unable to fetch tendermint txs for block:", processHeight, "ERR:", err)
				time.Sleep(sleepDuration)
				continue
			}

			err = handleTxSearchResult(initClientCtx, txSearchResult, &kafkaProducer, protoCodec)
			if err != nil {
				logging.Error("Unable to handle tendermint txs at height:", processHeight, "ERR:", err)
				time.Sleep(sleepDuration)
				continue
			}

			err = db.SetCosmosStatus(processHeight)
			if err != nil {
				logging.Fatal("ERROR setting tendermint status:", err)
			}

		}
		time.Sleep(sleepDuration)
	}
}
