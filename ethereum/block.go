package ethereum

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
	contracts2 "github.com/persistenceOne/persistenceBridge/ethereum/contracts"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func handleBlock(client *ethclient.Client, ctx *context.Context, block *types.Block, kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec) error {
	for _, transaction := range block.Transactions() {
		if transaction.To() != nil {
			var contract contracts2.ContractI
			switch transaction.To().String() {
			case contracts2.LiquidStaking.GetAddress():
				contract = &contracts2.LiquidStaking
			case contracts2.TokenWrapper.GetAddress():
				contract = &contracts2.TokenWrapper
			case contracts2.STokens.GetAddress():
				contract = &contracts2.STokens
			default:
			}
			if contract != nil {
				err := handleTransaction(client, ctx, transaction, contract, kafkaProducer, protoCodec)
				if err != nil {
					log.Printf("Failed to process ethereum tx: %s\n", transaction.Hash().String())
					return err
				}
			}
		}
	}
	return nil
}

func handleTransaction(client *ethclient.Client, ctx *context.Context, transaction *types.Transaction, contract contracts2.ContractI, kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec) error {
	receipt, err := client.TransactionReceipt(*ctx, transaction.Hash())
	if err != nil {
		log.Printf("Error while fetching receipt of tx %s: %s\n", transaction.Hash().String(), err.Error())
		return err
	}

	if receipt.Status == 1 {
		log.Printf("RECEIVED ETH Tx: %s\n", transaction.Hash().String())
		method, arguments, err := contract.GetMethodAndArguments(transaction.Data())
		if err != nil {
			log.Fatalf("Error in getting method and arguments of %s,: %s\n", contract.GetName(), err.Error())
			return err
		}

		if processFunc, ok := contract.GetMethods()[method.RawName]; ok {
			err = processFunc(kafkaProducer, protoCodec, arguments)
			if err != nil {
				log.Fatalf("Error in processing arguments of contarct %s method %s for tx %s: %s\n", contract.GetName(), method.RawName, transaction.Hash().String(), err.Error())
				return err
			}
		}
	}
	return nil
}
