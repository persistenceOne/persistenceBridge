package commands

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application"
	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	db2 "github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/rpc"
	"github.com/persistenceOne/persistenceBridge/application/shutdown"
	ethereum2 "github.com/persistenceOne/persistenceBridge/ethereum"
	"github.com/persistenceOne/persistenceBridge/kafka"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	tendermint2 "github.com/persistenceOne/persistenceBridge/tendermint"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func StartCommand() *cobra.Command {
	pBridgeCommand := &cobra.Command{
		Use:   "start",
		Short: "starts persistenceBridge",
		RunE: func(cmd *cobra.Command, args []string) error {

			homePath, err := cmd.Flags().GetString(constants2.FlagPBridgeHome)
			if err != nil {
				log.Fatalln(err)
			}

			showDebugLog, err := cmd.Flags().GetBool(constants2.FlagShowDebugLog)
			if err != nil {
				log.Fatalln(err)
			}
			logging.ShowDebugLog(showDebugLog)

			pStakeConfig := configuration.InitConfig()
			_, err = toml.DecodeFile(filepath.Join(homePath, "config.toml"), &pStakeConfig)
			if err != nil {
				log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
			}
			configuration.GetAppConfig().CASP.SetAPIToken()
			ethAddress, err := casp.GetEthAddress()
			if err != nil {
				log.Fatalln(err)
			}
			tmAddress, err := tendermint2.SetBech32PrefixesAndPStakeWrapAddress()
			if err != nil {
				log.Fatalln(err)
			}
			configuration.ValidateAndSeal()

			err = logging.InitializeBot()
			if err != nil {
				log.Fatalln(err)
			}

			logging.Info("Bridge (Wrap) Tendermint Address:", tmAddress.String())
			logging.Info("Bridge (Admin) Ethereum Address:", ethAddress.String())

			tmSleepTime, err := cmd.Flags().GetInt(constants2.FlagTendermintSleepTime)
			if err != nil {
				log.Fatalln(err)
			}

			tmStart, err := cmd.Flags().GetInt64(constants2.FlagTendermintStartHeight)
			if err != nil {
				log.Fatalln(err)
			}

			ethSleepTime, err := cmd.Flags().GetInt(constants2.FlagEthereumSleepTime)
			if err != nil {
				log.Fatalln(err)
			}

			ethStart, err := cmd.Flags().GetInt64(constants2.FlagEthereumStartHeight)
			if err != nil {
				log.Fatalln(err)
			}

			timeout, err := cmd.Flags().GetString(constants2.FlagTimeOut)
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

			chain, err := tendermint2.InitializeAndStartChain(timeout, homePath)
			if err != nil {
				log.Fatalln(err)
			}

			ethereumClient, err := ethclient.Dial(pStakeConfig.Ethereum.EthereumEndPoint)
			if err != nil {
				log.Fatalf("Error while dialing to eth orchestrator %s: %s\n", pStakeConfig.Ethereum.EthereumEndPoint, err.Error())
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
			kafkaState := utils.NewKafkaState(pStakeConfig.Kafka.Brokers, homePath, pStakeConfig.Kafka.TopicDetail)
			end := make(chan bool)
			ended := make(chan bool)
			go kafka.KafkaRoutine(kafkaState, protoCodec, chain, ethereumClient, end, ended)

			go rpc.StartServer(pStakeConfig.RPCEndpoint)

			logging.Info("Starting to listen ethereum....")
			go ethereum2.StartListening(ethereumClient, time.Duration(ethSleepTime)*time.Millisecond, pStakeConfig.Kafka.Brokers, protoCodec)

			logging.Info("Starting to listen tendermint....")
			go tendermint2.StartListening(clientContext, chain, pStakeConfig.Kafka.Brokers, protoCodec, time.Duration(tmSleepTime)*time.Millisecond)

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
					time.Sleep(100 * time.Millisecond)
				}
			}

			return nil
		},
	}
	pBridgeCommand.Flags().String(constants2.FlagTimeOut, constants2.DefaultTimeout, "timeout time for connecting to rpc")
	pBridgeCommand.Flags().String(constants2.FlagPBridgeHome, constants2.DefaultPBridgeHome, "home for pBridge")
	pBridgeCommand.Flags().Bool(constants2.FlagShowDebugLog, false, "show debug logs")
	pBridgeCommand.Flags().Int(constants2.FlagTendermintSleepTime, constants2.DefaultTendermintSleepTime, "sleep time between block checking for tendermint in ms")
	pBridgeCommand.Flags().Int(constants2.FlagEthereumSleepTime, constants2.DefaultEthereumSleepTime, "sleep time between block checking for ethereum in ms")
	pBridgeCommand.Flags().Int64(constants2.FlagTendermintStartHeight, constants2.DefaultTendermintStartHeight, fmt.Sprintf("Start checking height on tendermint chain from this height (default %d - starts from where last left)", constants2.DefaultTendermintStartHeight))
	pBridgeCommand.Flags().Int64(constants2.FlagEthereumStartHeight, constants2.DefaultEthereumStartHeight, fmt.Sprintf("Start checking height on ethereum chain from this height (default %d - starts from where last left)", constants2.DefaultEthereumStartHeight))

	return pBridgeCommand
}
