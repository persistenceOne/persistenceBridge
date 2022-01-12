/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package ethereum

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func onNewBlock(ctx context.Context, latestBlockHeight uint64, client *ethclient.Client, kafkaProducer sarama.SyncProducer, protoCodec *codec.ProtoCodec, database *badger.DB) error {
	return db.IterateOutgoingEthTx(database, func(key []byte, value []byte) error {
		var ethTx db.OutgoingEthereumTransaction

		err := json.Unmarshal(value, &ethTx)
		if err != nil {
			return fmt.Errorf("%w %s [ETH onNewBlock]: %s", ErrTxUnmarshal, string(key), err.Error())
		}

		txReceipt, err := client.TransactionReceipt(ctx, ethTx.TxHash)
		if err != nil {
			if txReceipt == nil && errors.Is(err, ethereum.NotFound) {
				logging.Info("Broadcast ethereum tx pending:", ethTx.TxHash)
			} else {
				logging.Error("Receipt fetch failed [onNewBlock] eth tx (need to check manually):", ethTx.TxHash.String(), "Error:", err)
			}
		}

		deleteTx := false

		if txReceipt.Status == 0 {
			logging.Error("Broadcast ethereum tx failed, Hash:", ethTx.TxHash.String(), "Block:", txReceipt.BlockNumber.Uint64())

			for _, msg := range ethTx.Messages {
				stakingAmount := sdkTypes.NewIntFromBigInt(msg.StakingAmount)
				wrapAmount := sdkTypes.NewIntFromBigInt(msg.WrapAmount)
				refundMsg := bankTypes.MsgSend{
					FromAddress: configuration.GetAppConfig().Tendermint.GetWrapAddress(),
					ToAddress:   msg.FromAddress.String(),
					Amount:      sdkTypes.NewCoins(sdkTypes.NewCoin(configuration.GetAppConfig().Tendermint.Denom, wrapAmount.Add(stakingAmount))),
				}

				msgBytes, err := protoCodec.MarshalInterface(&refundMsg)
				if err != nil {
					return fmt.Errorf("failed to generate msgBytes for tm refundMsg [ETH onNewBlock]: %w", err)
				}

				err = utils.ProducerDeliverMessage(msgBytes, utils.ToTendermint, kafkaProducer)
				if err != nil {
					logging.Error("Failed to add msg to kafka queue [ETH onNewBlock] ToEth, Message:", msg, "Error:", err)

					return err
				}
			}

			deleteTx = true
		} else {
			confirmedBlocks := latestBlockHeight - txReceipt.BlockNumber.Uint64()
			if confirmedBlocks >= constants.EthereumBlockConfirmations {
				logging.Info("Broadcast ethereum tx successful. Hash:", ethTx.TxHash, "Block:", txReceipt.BlockNumber.Uint64(), "Confirmed blocks:", confirmedBlocks)
				deleteTx = true
			} else {
				logging.Info("Broadcast ethereum tx confirmation pending. Hash:", ethTx.TxHash, "Block:", txReceipt.BlockNumber.Uint64(), "Confirmed blocks:", confirmedBlocks)
			}

			if deleteTx {
				return db.DeleteOutgoingEthereumTx(database, ethTx.TxHash)
			}

			return nil
		}

		return nil
	})
}
