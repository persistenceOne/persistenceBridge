package ethereum

import (
	"context"
	"github.com/cosmos/cosmos-sdk/codec"
	contracts2 "github.com/persistenceOne/persistenceBridge/ethereum/contracts"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func handleBlock(client *ethclient.Client, ctx *context.Context, block *types.Block, kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec) error {
	for _, transaction := range block.Transactions() {
		if transaction.To() != nil {
			var contract contracts2.ContractI
			switch transaction.To().String() {
			case contracts2.LiquidStaking.GetAddress():
				contract = &contracts2.LiquidStaking
			default:
			}
			if contract != nil {
				err := handleTransaction(client, ctx, transaction, contract, kafkaState, protoCodec)
				if err != nil {
					log.Printf("Failed to process ethereum tx: %s\n", transaction.Hash().String())
					return err
				}
			}
		}
	}
	return nil
}

func handleTransaction(client *ethclient.Client, ctx *context.Context, transaction *types.Transaction, contract contracts2.ContractI, kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec) error {
	receipt, err := client.TransactionReceipt(*ctx, transaction.Hash())
	if err != nil {
		log.Printf("Error while fetching receipt of tx %s: %s\n", transaction.Hash().String(), err.Error())
		return err
	}

	if receipt.Status == 1 {
		log.Printf("ETH Tx: %s\n", transaction.Hash().String())
		method, arguments, err := contract.GetMethodAndArguments(transaction.Data())
		if err != nil {
			log.Fatalf("Error in getting method and arguments of %s,: %s\n", contract.GetName(), err.Error())
			return err
		}

		if processFunc, ok := contract.GetMethods()[method.RawName]; ok {
			err = processFunc(kafkaState, protoCodec, arguments)
			if err != nil {
				log.Fatalf("Error in processing arguments of contarct %s method %s for tx %s: %s\n", contract.GetName(), method.RawName, transaction.Hash().String(), err.Error())
				return err
			}
		}
	}
	return nil
}
