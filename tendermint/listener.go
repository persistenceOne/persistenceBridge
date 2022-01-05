/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package tendermint

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/relayer/relayer"
	"github.com/dgraph-io/badger/v3"
	coreTypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/shutdown"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

const (
	finality  = 5
	bitlength = 64

	percent100       = 1
	blockPeriodRound = 10
)

func StartListening(initClientCtx *client.Context, database *badger.DB, chain *relayer.Chain, brokers []string, protoCodec *codec.ProtoCodec, sleepDuration time.Duration) {
	ctx := context.Background()
	kafkaProducer := utils.NewProducer(brokers, utils.SaramaConfig())

	var err error

	defer func(kafkaProducer sarama.SyncProducer) {
		err = kafkaProducer.Close()
		if err != nil {
			logging.Error("an error while closing Kafka", err)
		}
	}(kafkaProducer)

	slashingParamsResponse, err := QuerySlashingParams(chain)
	if err != nil {
		logging.Error("Params not found", "ERR:", err)
	}

	minSignedPerWindow, err := strconv.ParseFloat(slashingParamsResponse.Params.MinSignedPerWindow.String(), bitlength)
	if err != nil {
		logging.Error("Cannot convert MinSignedPerWindow to float, ERR:", err)
	}

	checkValidatorStatusPeriod := int64(float64(slashingParamsResponse.Params.SignedBlocksWindow) * (percent100 - minSignedPerWindow) / blockPeriodRound)

	for {
		// For Tendermint, we can directly query without waiting for blocks since there is finality
		err := onNewBlock(ctx, initClientCtx, chain, kafkaProducer, database, protoCodec)
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

			const timeToStopBridge = 5 * time.Second

			time.Sleep(timeToStopBridge)

			continue
		}

		abciInfo, err := chain.Client.ABCIInfo(ctx)
		if err != nil {
			logging.Error("Unable to fetch tendermint ABCI info:", err)

			time.Sleep(sleepDuration)

			continue
		}

		cosmosStatus, err := db.GetCosmosStatus(database)
		if err != nil {
			logging.Error("Stopping Tendermint Listener, unable to get status, Error:", err)

			shutdown.SetTMStopped(true)

			return
		}

		if (abciInfo.Response.LastBlockHeight - cosmosStatus.LastCheckHeight) > finality {
			processHeight := cosmosStatus.LastCheckHeight + 1

			logging.Info("Tendermint Block:", processHeight)

			resultTxs, err := getAllTxResults(ctx, chain, processHeight)
			if err != nil {
				time.Sleep(sleepDuration)

				continue
			}

			if processHeight%checkValidatorStatusPeriod == 0 {
				CheckValidators(chain, database, processHeight)
			}

			err = handleTxSearchResult(initClientCtx, resultTxs, kafkaProducer, database, protoCodec)
			if err != nil {
				logging.Error("Unable to handle tendermint txs at height:", processHeight, "ERR:", err)

				time.Sleep(sleepDuration)

				continue
			}

			err = db.SetCosmosStatus(database, processHeight)
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

	fmt.Println("!!!!!!!!!", chain == nil)
	fmt.Println("!!!!!!!!!1", chain.Client == nil)

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
