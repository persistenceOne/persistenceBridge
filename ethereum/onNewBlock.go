package ethereum

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
	"github.com/pkg/errors"
	"math/big"
)

func onNewBlock(ctx context.Context, latestBlockHeight uint64, client *ethclient.Client, kafkaProducer *sarama.SyncProducer) error {
	return db.IterateOutgoingEthTx(func(key []byte, value []byte) error {
		var ethTx db.OutgoingEthereumTransaction
		err := json.Unmarshal(value, &ethTx)
		if err != nil {
			return fmt.Errorf("failed to unmarshal OutgoingEthereumTransaction %s [ETH onNewBlock]: %s", string(key), err.Error())
		}
		fmt.Println(ethTx.TxHash)
		transaction, _, err := client.TransactionByHash(ctx, ethTx.TxHash)
		txReceipt, err := client.TransactionReceipt(ctx, ethTx.TxHash)
		fmt.Println(txReceipt.Logs)
		sender, err := client.TransactionSender(ctx, transaction, txReceipt.BlockHash,txReceipt.TransactionIndex)
		fmt.Println(sender.Hash())
		if err!= nil {
			fmt.Println(err.Error())
		}
		if err != nil {
			if txReceipt == nil && err == ethereum.NotFound {
				logging.Info("Broadcast ethereum tx pending:", ethTx.TxHash)
			} else {
				logging.Error("Receipt fetch failed [onNewBlock] eth tx:", ethTx.TxHash.String(), "Error:", err)
			}
		} else {
			deleteTx := false
			if txReceipt.Status == 0 {
				err := getError(ctx, sender, transaction, txReceipt.BlockNumber, *client)
				if err!=nil {
					return err
				}
				logging.Error("Broadcast ethereum tx failed, Hash:", ethTx.TxHash.String(), "Block:", txReceipt.BlockNumber.Uint64())
				for _, msg := range ethTx.Messages {
					msgBytes, err := json.Marshal(msg)
					if err != nil {
						return err
					}
					err = utils.ProducerDeliverMessage(msgBytes, utils.ToEth, *kafkaProducer)
					if err != nil {
						logging.Error("Failed to add msg to kafka queue [ETH onNewBlock] ToEth, Message:", msg, "Error:", err)
						return err
					}
				}
				deleteTx = true
			} else {
				confirmedBlocks := latestBlockHeight - txReceipt.BlockNumber.Uint64()
				if confirmedBlocks >= 12 {
					logging.Info("Broadcast ethereum tx successful. Hash:", ethTx.TxHash, "Block:", txReceipt.BlockNumber.Uint64(), "Confirmed blocks:", confirmedBlocks)
					deleteTx = true
				} else {
					logging.Info("Broadcast ethereum tx confirmation pending. Hash:", ethTx.TxHash, "Block:", txReceipt.BlockNumber.Uint64(), "Confirmed blocks:", confirmedBlocks)
				}
			}
			if deleteTx {
				return db.DeleteOutgoingEthereumTx(ethTx.TxHash)
			}
			return nil
		}
		return nil
	})


}


func getError(ctx context.Context, from common.Address, tx *types.Transaction, blockNum *big.Int, ethClient ethclient.Client) error {
	msg := ethereum.CallMsg{
		From:     from,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}

	res, err := ethClient.CallContract(ctx, msg, blockNum)
	if err != nil {
		return errors.Wrap(err, "Cannot get revert reason fatal")
	}
	if len(res) == 0 {
		return errors.New("Out of gas")
	}

	return unpackError(res)
}

var (
	errorSig            = []byte{0x08, 0xc3, 0x79, 0xa0} // Keccak256("Error(string)")[:4]
	abiString, _        = abi.NewType("string", "", nil)
)

func unpackError(result []byte) error {
	if !bytes.Equal(result[:4], errorSig) {
		return errors.New("TX result not of type Error(string)")
	}
	vs, err := abi.Arguments{{Type: abiString}}.UnpackValues(result[4:])
	if err != nil {
		return errors.Wrap(err, "Unpacking revert reason")
	}
	return errors.New(vs[0].(string))
}