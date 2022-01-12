/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package tendermint

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dgraph-io/badger/v3"
	goEthCommon "github.com/ethereum/go-ethereum/common"
	tmCoreTypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func handleTxSearchResult(clientCtx *client.Context, resultTxs []*tmCoreTypes.ResultTx, kafkaProducer sarama.SyncProducer, database *badger.DB, protoCodec *codec.ProtoCodec) error {
	for i, transaction := range resultTxs {
		logging.Info("Tendermint TX:", transaction.Hash.String(), fmt.Sprintf("(%d)", i+1))

		err := collectAllWrapAndRevertTxs(clientCtx, database, transaction)
		if err != nil {
			logging.Error("Failed to process tendermint transaction:", transaction.Hash.String())

			return err
		}
	}

	wrapOrRevert(kafkaProducer, protoCodec, database)

	return nil
}

func collectAllWrapAndRevertTxs(clientCtx *client.Context, database *badger.DB, txQueryResult *tmCoreTypes.ResultTx) error {
	if txQueryResult.TxResult.GetCode() == 0 {
		// Should be used if txQueryResult.Tx is string: `decodedTx, err := base64.StdEncoding.DecodeString(txQueryResult.Tx)`
		txInterface, err := clientCtx.TxConfig.TxDecoder()(txQueryResult.Tx)
		if err != nil {
			return err
		}

		transaction, ok := txInterface.(signing.Tx)
		if !ok {
			return err
		}

		memo := strings.TrimSpace(transaction.GetMemo())

		for i, msg := range transaction.GetMsgs() {
			switch txMsg := msg.(type) {
			case *bankTypes.MsgSend:
				if txMsg.ToAddress == configuration.GetAppConfig().Tendermint.GetWrapAddress() {
					if memo != "DO_NOT_REVERT" {
						for _, coin := range txMsg.Amount {
							// Do not check for TendermintTxToKafka exists.
							if !db.CheckIncomingTendermintTxExists(database, txQueryResult.Hash, uint(i), coin.Denom) {
								err = db.AddIncomingTendermintTx(database, &db.IncomingTendermintTx{
									TxHash:      txQueryResult.Hash,
									MsgIndex:    uint(i),
									Denom:       coin.Denom,
									FromAddress: txMsg.FromAddress,
									Amount:      coin.Amount,
									Memo:        memo,
								})
								if err != nil {
									return err
								}

								err = db.AddTendermintTxToKafka(database, db.TendermintTxToKafka{
									TxHash:   txQueryResult.Hash,
									MsgIndex: uint(i),
									Denom:    coin.Denom,
								})
								if err != nil {
									// nolint already has %w
									// nolint: errorlint
									return fmt.Errorf("%s: %v", ErrPartialSend, err)
								}
							}
						}
					} else {
						logging.Info("Deposited to wrapping address, TxHash:", txQueryResult.Hash.String(), "amount:", txMsg.Amount.String())
					}
				}
			default:
			}
		}
	}

	return nil
}

func wrapOrRevert(kafkaProducer sarama.SyncProducer, protoCodec *codec.ProtoCodec, database *badger.DB) {
	tmTxToKafkaList, err := db.GetAllTendermintTxToKafka(database)
	if err != nil {
		logging.Fatal(err)
	}

	for _, tmTxToKafka := range tmTxToKafkaList {
		tx, err := db.GetIncomingTendermintTx(database, tmTxToKafka.TxHash, tmTxToKafka.MsgIndex, tmTxToKafka.Denom)
		if err != nil {
			logging.Fatalf("%s [TM Listener]: %s", ErrGetIncomingTendermintTx, err.Error())
		}

		validMemo := true

		var (
			ethAddress goEthCommon.Address
			ratio      sdk.Dec
		)

		ethAddress, ratio, err = getWrapAddressAndStakingRatio(tx.Memo)
		if err != nil {
			validMemo = false
		}

		if tx.Denom == configuration.GetAppConfig().Tendermint.Denom && validMemo && tx.Amount.GTE(sdk.NewInt(configuration.GetAppConfig().Tendermint.MinimumWrapAmount)) {
			logging.Info("Tendermint Wrap Tx:", tx.TxHash, "Msg Index:", tx.MsgIndex, "Amount:", tx.Amount.String())

			fromAddress, err := sdk.AccAddressFromBech32(tx.FromAddress)
			if err != nil {
				logging.Fatal(err)
			}

			stakingAmount := sdk.NewDecFromInt(tx.Amount).Mul(ratio).TruncateInt()
			wrapAmount := tx.Amount.Sub(stakingAmount)

			ethTxMsg := db.NewWrapTokenMsg(fromAddress, tx.TxHash, stakingAmount.BigInt(), ethAddress, wrapAmount.BigInt())

			var msgBytes []byte

			msgBytes, err = json.Marshal(ethTxMsg)
			if err != nil {
				logging.Fatal(err)
			}

			logging.Info("DIRECT STAKING: Adding wrap & stake token msg to kafka producer ToEth, from:", tx.FromAddress, "to:", ethAddress.String(), "wrap amount:", wrapAmount.String(), "staking amount:", stakingAmount.String())

			err = utils.ProducerDeliverMessage(msgBytes, utils.ToEth, kafkaProducer)
			if err != nil {
				logging.Fatalf("%s ToEth [TM Listener]: %s", ErrAddToKafkaQueue, err.Error())
			}

			if stakingAmount.GT(sdk.ZeroInt()) {
				stakingMsg := &stakingTypes.MsgDelegate{
					DelegatorAddress: configuration.GetAppConfig().Tendermint.GetWrapAddress(),
					ValidatorAddress: "",
					Amount:           sdk.NewCoin(configuration.GetAppConfig().Tendermint.Denom, stakingAmount),
				}
				stakingMsgBytes, err := protoCodec.MarshalInterface(stakingMsg)
				if err != nil {
					logging.Fatal("Tendermint Listener: Staking Message marshalling failed:", err.Error())
				}
				logging.Info("DIRECT STAKING: Adding direct staking msg to kafka producer MsgDelegate, from:", tx.FromAddress, "to:", ethAddress.String(), "direct staking amount:", stakingAmount.String())
				err = utils.ProducerDeliverMessage(stakingMsgBytes, utils.MsgDelegate, kafkaProducer)
				if err != nil {
					logging.Fatalf("failed to add staking msg to kafka queue MsgDelegate [TM Listener]: %s", err.Error())
				}
			}
		} else {
			revertToken := sdk.NewCoin(tx.Denom, tx.Amount)

			logging.Info("Reverting coins,TxHash:", tx.TxHash, "Msg Index:", tx.MsgIndex, "Coin:", revertToken.String())

			revertCoins(tx.FromAddress, sdk.NewCoins(revertToken), kafkaProducer, protoCodec)
		}

		err = db.DeleteTendermintTxToKafka(database, tx.TxHash, tx.MsgIndex, tx.Denom)
		if err != nil {
			logging.Fatal(err)
		}
	}
}

func revertCoins(toAddress string, coins sdk.Coins, kafkaProducer sarama.SyncProducer, protoCodec *codec.ProtoCodec) {
	msg := &bankTypes.MsgSend{
		FromAddress: configuration.GetAppConfig().Tendermint.GetWrapAddress(),
		ToAddress:   toAddress,
		Amount:      coins,
	}

	msgBytes, err := protoCodec.MarshalInterface(msg)
	if err != nil {
		logging.Fatal(err)
	}

	logging.Info("REVERT COINS: adding send coin msg to kafka producer MsgSend, to:", toAddress, "amount:", coins.String())

	err = utils.ProducerDeliverMessage(msgBytes, utils.MsgSend, kafkaProducer)
	if err != nil {
		logging.Fatalf("%s ToEth [TM Listener REVERT]: %s", ErrAddToKafkaQueue, err.Error())
	}
}

func getWrapAddressAndStakingRatio(memo string) (goEthCommon.Address, sdk.Dec, error) {
	split := strings.Split(memo, ",")
	if len(split) > 2 || len(split) == 0 {
		return goEthCommon.Address{}, sdk.Dec{}, fmt.Errorf("invalid memo for bridge")
	}

	isEthAddress := false

	var ethAddress goEthCommon.Address
	if goEthCommon.IsHexAddress(split[0]) {
		ethAddress = goEthCommon.HexToAddress(split[0])
		isEthAddress = ethAddress != constants.EthereumZeroAddress()
	}

	if !isEthAddress {
		return goEthCommon.Address{}, sdk.Dec{}, fmt.Errorf("invalid memo for bridge")
	}

	if len(split) == 1 {
		return ethAddress, sdk.ZeroDec(), nil
	}

	ratio, err := sdk.NewDecFromStr(split[1])

	if err != nil {
		return goEthCommon.Address{}, sdk.Dec{}, fmt.Errorf("invalid memo for bridge")
	}

	if ratio.IsNegative() {
		return goEthCommon.Address{}, sdk.Dec{}, fmt.Errorf("negative ratio: invalid memo for bridge")
	}

	return ethAddress, ratio, nil
}
