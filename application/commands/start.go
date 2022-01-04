/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application"
	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	db2 "github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/rpc"
	"github.com/persistenceOne/persistenceBridge/application/shutdown"
	"github.com/persistenceOne/persistenceBridge/ethereum"
	"github.com/persistenceOne/persistenceBridge/kafka"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/tendermint"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartCommand() *cobra.Command {
	pBridgeCommand := &cobra.Command{
		Use:   "start",
		Short: "starts persistenceBridge",
		RunE: func(cmd *cobra.Command, args []string) error {

			homePath, err := cmd.Flags().GetString(constants.FlagPBridgeHome)
			if err != nil {
				log.Fatalln(err)
			}

			showDebugLog, err := cmd.Flags().GetBool(constants.FlagShowDebugLog)
			if err != nil {
				log.Fatalln(err)
			}
			logging.ShowDebugLog(showDebugLog)

			err = logging.InitializeBot()
			if err != nil {
				log.Fatalln(err)
			}

			setAndSealConfig(homePath)

			tmSleepTime, err := cmd.Flags().GetInt(constants.FlagTendermintSleepTime)
			if err != nil {
				log.Fatalln(err)
			}

			tmStart, err := cmd.Flags().GetInt64(constants.FlagTendermintStartHeight)
			if err != nil {
				log.Fatalln(err)
			}

			ethSleepTime, err := cmd.Flags().GetInt(constants.FlagEthereumSleepTime)
			if err != nil {
				log.Fatalln(err)
			}

			ethStart, err := cmd.Flags().GetInt64(constants.FlagEthereumStartHeight)
			if err != nil {
				log.Fatalln(err)
			}

			timeout, err := cmd.Flags().GetString(constants.FlagTimeOut)
			if err != nil {
				log.Fatalln(err)
			}

			db, err := db2.InitializeDB(homePath+"/db", tmStart, ethStart)
			if err != nil {
				log.Fatalln(err)
			}
			defer func(db *badger.DB) {
				err := db.Close()
				if err != nil {
					log.Println("Error while closing DB: ", err.Error())
				}
			}(db)

			unboundEpochTime, err := db2.GetUnboundEpochTime()
			if err != nil {
				log.Fatalln(err)
			}
			log.Printf("unbound epoch time: %d\n", unboundEpochTime.Epoch)

			validators, err := db2.GetValidators()
			if err != nil {
				log.Fatalln(err)
			}
			if len(validators) == 0 {
				log.Fatalln("no validator has been added")
			} else {
				for i, validator := range validators {
					fmt.Println(fmt.Sprintf("%d. Name: %s, Address: %s", i+1, validator.Name, validator.Address))
				}
			}

			chain, err := tendermint.InitializeAndStartChain(timeout, homePath)
			if err != nil {
				log.Fatalln(err)
			}

			ethereumClient, err := ethclient.Dial(configuration.GetAppConfig().Ethereum.EthereumEndPoint)
			if err != nil {
				log.Fatalf("Error while dialing to eth orchestrator %s: %s\n", configuration.GetAppConfig().Ethereum.EthereumEndPoint, err.Error())
			}

			encodingConfig := application.MakeEncodingConfig()
			clientContext := client.Context{}.
				WithJSONMarshaler(encodingConfig.Marshaler).
				WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
				WithTxConfig(encodingConfig.TransactionConfig).
				WithLegacyAmino(encodingConfig.Amino).
				WithInput(os.Stdin).
				WithAccountRetriever(authTypes.AccountRetriever{}).
				WithBroadcastMode(configuration.GetAppConfig().Tendermint.BroadcastMode).
				WithHomeDir(homePath)

			protoCodec := codec.NewProtoCodec(clientContext.InterfaceRegistry)
			kafkaState := utils.NewKafkaState(configuration.GetAppConfig().Kafka.Brokers, homePath, configuration.GetAppConfig().Kafka.TopicDetail)
			end := make(chan bool)
			ended := make(chan bool)
			go kafka.KafkaRoutine(kafkaState, protoCodec, chain, ethereumClient, end, ended)

			go rpc.StartServer(configuration.GetAppConfig().RPCEndpoint)

			logging.Info("Starting to listen ethereum....")
			go ethereum.StartListening(ethereumClient, time.Duration(ethSleepTime)*time.Millisecond, configuration.GetAppConfig().Kafka.Brokers, protoCodec)

			logging.Info("Starting to listen tendermint....")
			go tendermint.StartListening(clientContext, chain, configuration.GetAppConfig().Kafka.Brokers, protoCodec, time.Duration(tmSleepTime)*time.Millisecond)

			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
			for sig := range signalChan {
				logging.Info("STOP SIGNAL RECEIVED:", sig.String(), "(Might take around a minute to stop)")
				shutdown.StopBridge()
				for {
					if !shutdown.GetKafkaConsumerClosed() {
						logging.Info("Stopping Kafka Routine!!!")
						kafka.KafkaClose(kafkaState, end, ended)()
						shutdown.SetKafkaConsumerClosed(true)
					}
					if shutdown.GetTMStopped() && shutdown.GetETHStopped() && shutdown.GetKafkaConsumerClosed() {
						return nil
					}
					time.Sleep(100 * time.Millisecond) // thread is put to sleep to prevent 100% CPU usage
				}
			}

			return nil
		},
	}
	pBridgeCommand.Flags().String(constants.FlagTimeOut, constants.DefaultTimeout, "timeout time for connecting to rpc")
	pBridgeCommand.Flags().String(constants.FlagPBridgeHome, constants.DefaultPBridgeHome, "home for pBridge")
	pBridgeCommand.Flags().Bool(constants.FlagShowDebugLog, false, "show debug logs")
	pBridgeCommand.Flags().Int(constants.FlagTendermintSleepTime, constants.DefaultTendermintSleepTime, "sleep time between block checking for tendermint in ms")
	pBridgeCommand.Flags().Int(constants.FlagEthereumSleepTime, constants.DefaultEthereumSleepTime, "sleep time between block checking for ethereum in ms")
	pBridgeCommand.Flags().Int64(constants.FlagTendermintStartHeight, constants.DefaultTendermintStartHeight, fmt.Sprintf("Start checking height on tendermint chain from this height (default %d - starts from where last left)", constants.DefaultTendermintStartHeight))
	pBridgeCommand.Flags().Int64(constants.FlagEthereumStartHeight, constants.DefaultEthereumStartHeight, fmt.Sprintf("Start checking height on ethereum chain from this height (default %d - starts from where last left)", constants.DefaultEthereumStartHeight))

	return pBridgeCommand
}

func setAndSealConfig(homePath string) {
	configuration.InitializeConfigFromFile(homePath)
	ethAddress, err := casp.GetEthAddress()
	if err != nil {
		log.Fatalln(err)
	}
	tmAddress, err := casp.GetTendermintAddress()
	if err != nil {
		log.Fatalln(err)
	}
	configuration.SetCASPAddresses(tmAddress, ethAddress)

	logging.Info("Bridge (Wrap) Tendermint Address:", configuration.GetAppConfig().Tendermint.GetWrapAddress())
	logging.Info("Bridge (Admin) Ethereum Address:", configuration.GetAppConfig().Ethereum.GetBridgeAdminAddress().String())

	configuration.ValidateAndSeal()
}
