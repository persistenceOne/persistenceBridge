package tendermint

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/relayer/relayer"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/shutdown"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
	coreTypes "github.com/tendermint/tendermint/rpc/core/types"
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
			shutdown.SetTMStopped(true)
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
			logging.Error("Stopping Tendermint Listener, unable to get status, Error:", err)
			shutdown.SetTMStopped(true)
			return
		}

		if abciInfo.Response.LastBlockHeight > cosmosStatus.LastCheckHeight {
			processHeight := cosmosStatus.LastCheckHeight + 1
			logging.Info("Tendermint Block:", processHeight)

			resultTxs, err := getAllTxResults(ctx, chain, processHeight)
			if err != nil {
				time.Sleep(sleepDuration)
				continue
			}

			validatorList, err := db.GetValidators()
			err = handleSlashedOrAboutToBeSlashed(chain, validatorList, processHeight)
			if err != nil {
				logging.Error("Unable to handle jailed or to be jailed check at height:", processHeight, "ERR:", err)
				continue
			}

			err = handleTxSearchResult(initClientCtx, resultTxs, &kafkaProducer, protoCodec)
			if err != nil {
				logging.Error("Unable to handle tendermint txs at height:", processHeight, "ERR:", err)
				time.Sleep(sleepDuration)
				continue
			}

			err = db.SetCosmosStatus(processHeight)
			if err != nil {
				logging.Error("Stopping Tendermint Listener, unable to set (DB) status to", processHeight, "Error:", err)
				shutdown.SetTMStopped(true)
				return
			}

		}
		time.Sleep(sleepDuration)
	}
}

func getAllTxResults(ctx context.Context, chain *relayer.Chain, height int64) ([]*coreTypes.ResultTx, error) {
	var resultTxs []*coreTypes.ResultTx
	page := 1
	txsMaxPerPage := 100
	txSearchResult, err := chain.Client.TxSearch(ctx, fmt.Sprintf("tx.height=%d", height), true, &page, &txsMaxPerPage, "asc")
	if err != nil {
		logging.Error("Unable to fetch tendermint txs for block:", height, "page:", page, "ERR:", err)
		return resultTxs, err
	}
	if txSearchResult.TotalCount <= txsMaxPerPage {
		return txSearchResult.Txs, nil
	}
	resultTxs = append(resultTxs, txSearchResult.Txs...)
	totalPages := int(math.Ceil(float64(txSearchResult.TotalCount) / float64(txsMaxPerPage)))
	for i := page + 1; i <= totalPages; i++ {
		txSearchResult, err = chain.Client.TxSearch(ctx, fmt.Sprintf("tx.height=%d", height), true, &i, &txsMaxPerPage, "asc")
		if err != nil {
			logging.Error("Unable to fetch tendermint txs for block:", height, "page:", i, "ERR:", err)
			return resultTxs, err
		}
		resultTxs = append(resultTxs, txSearchResult.Txs...)
	}
	return resultTxs, nil
}
