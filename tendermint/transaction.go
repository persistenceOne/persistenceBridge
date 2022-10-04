/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package tendermint

import (
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	goEthCommon "github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
	tmCoreTypes "github.com/tendermint/tendermint/rpc/core/types"
)

func handleTxSearchResult(clientCtx client.Context, resultTxs []*tmCoreTypes.ResultTx, kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec) error {
	for i, transaction := range resultTxs {
		logging.Info("Tendermint TX:", transaction.Hash.String(), fmt.Sprintf("(%d)", i+1))
		err := collectAllWrapAndRevertTxs(clientCtx, transaction)
		if err != nil {
			logging.Error("Failed to process tendermint transaction:", transaction.Hash.String())
			return err
		}
	}
	wrapOrRevert(kafkaProducer, protoCodec)
	return nil
}

func collectAllWrapAndRevertTxs(clientCtx client.Context, txQueryResult *tmCoreTypes.ResultTx) error {
	if txQueryResult.TxResult.GetCode() == 0 {
		// Should be used if txQueryResult.Tx is string
		//decodedTx, err := base64.StdEncoding.DecodeString(txQueryResult.Tx)
		//if err != nil {
		//	return txMsgs, err
		//}

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
				if txMsg.ToAddress == configuration.GetAppConfig().Tendermint.GetWrapAddress() && txMsg.ToAddress != txMsg.FromAddress {
					if memo != "DO_NOT_REVERT" {
						for _, coin := range txMsg.Amount {
							// Do not check for TendermintTxToKafka exists.
							if !db.CheckIncomingTendermintTxExists(txQueryResult.Hash, uint(i), coin.Denom) {
								err = db.AddIncomingTendermintTx(db.IncomingTendermintTx{
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
								err = db.AddTendermintTxToKafka(db.TendermintTxToKafka{
									TxHash:   txQueryResult.Hash,
									MsgIndex: uint(i),
									Denom:    coin.Denom,
								})
								if err != nil {
									return fmt.Errorf("added to IncomingTendermintTx but NOT to TendermintTxToKafka failed. Tx won't be added to kafka: %v", err)
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

func wrapOrRevert(kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec) {
	tmTxToKafkaList, err := db.GetAllTendermintTxToKafka()
	if err != nil {
		logging.Fatal(err)
	}
	for _, tmTxToKafka := range tmTxToKafkaList {
		tx, err := db.GetIncomingTendermintTx(tmTxToKafka.TxHash, tmTxToKafka.MsgIndex, tmTxToKafka.Denom)
		if err != nil {
			logging.Fatal(fmt.Errorf("failed to get IncomingTendermintTx by TendermintTxToKafka [TM Listener]: %s", err.Error()))
		}
		revertToken := sdk.NewCoin(tx.Denom, tx.Amount)
		logging.Info("Reverting coins,TxHash:", tx.TxHash, "Msg Index:", tx.MsgIndex, "Coin:", revertToken.String())
		revertCoins(tx.FromAddress, sdk.NewCoins(revertToken), kafkaProducer, protoCodec)
		err = db.DeleteTendermintTxToKafka(tx.TxHash, tx.MsgIndex, tx.Denom)
		if err != nil {
			logging.Fatal(err)
		}
	}
}

func revertCoins(toAddress string, coins sdk.Coins, kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec) {
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
	err = utils.ProducerDeliverMessage(msgBytes, utils.MsgSend, *kafkaProducer)
	if err != nil {
		logging.Fatal(fmt.Errorf("failed to add msg to kafka queue ToEth [TM Listener REVERT]: %s", err.Error()))
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
		isEthAddress = ethAddress.String() != constants.EthereumZeroAddress
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
	if ratio.GT(sdk.NewDec(1)) {
		return goEthCommon.Address{}, sdk.Dec{}, fmt.Errorf("invalid memo for bridge, ratio cannot be greater than 1")
	}
	return ethAddress, ratio, nil
}
