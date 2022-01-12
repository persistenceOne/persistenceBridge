/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"

	"github.com/persistenceOne/persistenceBridge/application"
	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/rpc"
	"github.com/persistenceOne/persistenceBridge/application/shutdown"
	"github.com/persistenceOne/persistenceBridge/ethereum"
	"github.com/persistenceOne/persistenceBridge/kafka"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/tendermint"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
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

			pStakeConfig := configuration.GetAppConfig()
			_, err = toml.DecodeFile(filepath.Join(homePath, "config.toml"), &pStakeConfig)
			if err != nil {
				log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
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

			// fixme: do not change global state
			setAndSealConfig(homePath)

			// refresh config after change
			pStakeConfig = configuration.GetAppConfig()

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

			database, err := db.InitializeDB(homePath+"/db", tmStart, ethStart)
			if err != nil {
				log.Fatalln(err)
			}

			defer func(database *badger.DB) {
				innerErr := database.Close()
				if innerErr != nil {
					log.Println("Error while closing DB: ", innerErr.Error())
				}
			}(database)

			unboundEpochTime, err := db.GetUnboundEpochTime(database)
			if err != nil {
				log.Fatalln(err)
			}

			log.Printf("unbound epoch time: %d\n", unboundEpochTime.Epoch)

			validators, err := db.GetValidators(database)
			if err != nil {
				log.Fatalln(err)
			}

			if len(validators) == 0 {
				log.Fatalln("no validator has been added")
			} else {
				for i, validator := range validators {
					fmt.Printf("%d. Name: %s, Address: %s\n", i+1, validator.Name, validator.Address)
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
			kafkaState := utils.NewKafkaState(pStakeConfig.Kafka.Brokers, homePath, constants.ToKafkaTopicDetail(pStakeConfig.Kafka.TopicDetails))

			// fixme: should they be bufferized chan struct{}?
			end := make(chan bool)
			ended := make(chan bool)

			go kafka.Routine(kafkaState, database, protoCodec, chain, ethereumClient, end, ended)

			go rpc.StartServer(pStakeConfig.RPCEndpoint, database)

			logging.Info("Starting to listen ethereum....")
			go ethereum.StartListening(ethereumClient, database, time.Duration(ethSleepTime)*time.Millisecond, pStakeConfig.Kafka.Brokers, protoCodec)

			logging.Info("Starting to listen tendermint....")
			go tendermint.StartListening(&clientContext, database, chain, pStakeConfig.Kafka.Brokers, protoCodec, time.Duration(tmSleepTime)*time.Millisecond)

			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

			const stopCheckPeriod = 100 * time.Millisecond

			for sig := range signalChan {
				logging.Info("STOP SIGNAL RECEIVED:", sig.String(), "(Might take around a minute to stop)")

				shutdown.StopBridge()

				for {
					if !shutdown.GetKafkaConsumerClosed() {
						logging.Info("Stopping Kafka Routine!!!")
						kafka.Close(kafkaState, end, ended)()
						shutdown.SetKafkaConsumerClosed(true)
					}

					if shutdown.GetTMStopped() && shutdown.GetETHStopped() && shutdown.GetKafkaConsumerClosed() {
						return nil
					}

					time.Sleep(stopCheckPeriod)
				}
			}

			return nil
		},
	}

	pBridgeCommand.Flags().String(constants.FlagTimeOut, constants.DefaultTimeout, "timeout time for connecting to rpc")
	pBridgeCommand.Flags().String(constants.FlagPBridgeHome, constants.DefaultPBridgeHome(), "home for pBridge")
	pBridgeCommand.Flags().Bool(constants.FlagShowDebugLog, false, "show debug logs")
	pBridgeCommand.Flags().Int(constants.FlagTendermintSleepTime, constants.DefaultTendermintSleepTime, "sleep time between block checking for tendermint in ms")
	pBridgeCommand.Flags().Int(constants.FlagEthereumSleepTime, constants.DefaultEthereumSleepTime, "sleep time between block checking for ethereum in ms")
	pBridgeCommand.Flags().Int64(constants.FlagTendermintStartHeight, constants.DefaultTendermintStartHeight, fmt.Sprintf("Start checking height on tendermint chain from this height (default %d - starts from where last left)", constants.DefaultTendermintStartHeight))
	pBridgeCommand.Flags().Int64(constants.FlagEthereumStartHeight, constants.DefaultEthereumStartHeight, fmt.Sprintf("Start checking height on ethereum chain from this height (default %d - starts from where last left)", constants.DefaultEthereumStartHeight))

	return pBridgeCommand
}

func setAndSealConfig(homePath string) {
	pStakeConfig := configuration.GetAppConfig()
	_, err := toml.DecodeFile(filepath.Join(homePath, "config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}

	// fixme init context with proper dependancies and timeout
	ctx := context.Background()

	ethAddress, err := casp.GetEthAddress(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	tmAddress, err := tendermint.SetBech32PrefixesAndPStakeWrapAddress(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	configuration.SetCASPAddresses(tmAddress, ethAddress)

	logging.Info("Bridge (Wrap) Tendermint Address:", configuration.GetAppConfig().Tendermint.GetWrapAddress())
	logging.Info("Bridge (Admin) Ethereum Address:", configuration.GetAppConfig().Ethereum.GetBridgeAdminAddress().String())

	configuration.ValidateAndSeal()
}
