package commands

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	db2 "github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/shutdown"
	ethereum2 "github.com/persistenceOne/persistenceBridge/ethereum"
	"github.com/persistenceOne/persistenceBridge/kafka"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	tendermint2 "github.com/persistenceOne/persistenceBridge/tendermint"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func StartCommand(initClientCtx client.Context) *cobra.Command {
	pBridgeCommand := &cobra.Command{
		Use:   "start [path_to_chain_json]",
		Short: "Start persistenceBridge",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			homePath, err := cmd.Flags().GetString(constants2.FlagPBridgeHome)
			if err != nil {
				log.Fatalln(err)
			}

			pStakeConfig := configuration.Config{}
			_, err = toml.DecodeFile(filepath.Join(homePath, "config.toml"), &pStakeConfig)
			if err != nil {
				log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
			}
			pStakeConfig = UpdateConfig(cmd, pStakeConfig)
			configuration.SetAppConfig(pStakeConfig)

			tmSleepTime, err := cmd.Flags().GetInt(constants2.FlagTendermintSleepTime)
			if err != nil {
				log.Fatalln(err)
			}

			coinType, err := cmd.Flags().GetUint32(constants2.FlagCoinType)
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

			ethereumEndPoint, err := cmd.Flags().GetString(constants2.FlagEthereumEndPoint)
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

			chain, err := tendermint2.InitializeAndStartChain(args[0], timeout, homePath, coinType)
			if err != nil {
				log.Fatalln(err)
			}

			ethereumClient, err := ethclient.Dial(ethereumEndPoint)
			if err != nil {
				log.Fatalf("Error while dialing to eth orchestrator %s: %s\n", ethereumEndPoint, err.Error())
			}

			protoCodec := codec.NewProtoCodec(initClientCtx.InterfaceRegistry)
			kafkaState := utils.NewKafkaState(pStakeConfig.Kafka.Brokers, homePath, pStakeConfig.Kafka.TopicDetail)
			go kafka.KafkaRoutine(kafkaState, protoCodec, chain, ethereumClient)

			log.Println("Starting to listen ethereum....")
			go ethereum2.StartListening(ethereumClient, time.Duration(ethSleepTime)*time.Millisecond, kafkaState, protoCodec)

			log.Println("Starting to listen tendermint....")
			go tendermint2.StartListening(initClientCtx.WithHomeDir(homePath), chain, kafkaState, protoCodec, time.Duration(tmSleepTime)*time.Millisecond)

			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
			for sig := range signalChan {
				log.Println("signal received to close: " + sig.String())
				shutdown.SetBridgeStopSignal(true)
				kafkaClosed := false
				for {
					if !kafkaClosed {
						log.Println("Stopping Kafka Routine!!!")
						kafka.KafkaClose(kafkaState)
						kafkaClosed = true
					}
					if shutdown.GetTMStopped() && shutdown.GetETHStopped() {
						return nil
					}
				}
			}

			return nil
		},
	}
	pBridgeCommand.Flags().String(constants2.FlagTimeOut, constants2.DefaultTimeout, "timeout time for connecting to rpc")
	pBridgeCommand.Flags().Uint32(constants2.FlagCoinType, constants2.DefaultCoinType, "coin type for wallet")
	pBridgeCommand.Flags().String(constants2.FlagPBridgeHome, constants2.DefaultPBridgeHome, "home for pBridge")
	pBridgeCommand.Flags().String(constants2.FlagEthereumEndPoint, constants2.DefaultEthereumEndPoint, "ethereum orchestrator to connect")
	pBridgeCommand.Flags().String(constants2.FlagKafkaPorts, constants2.DefaultKafkaPorts, "ports kafka brokers are running on, --ports 192.100.10.10:443,192.100.10.11:443")
	pBridgeCommand.Flags().Int(constants2.FlagTendermintSleepTime, constants2.DefaultTendermintSleepTime, "sleep time between block checking for tendermint in ms")
	pBridgeCommand.Flags().Int(constants2.FlagEthereumSleepTime, constants2.DefaultEthereumSleepTime, "sleep time between block checking for ethereum in ms")
	pBridgeCommand.Flags().Int64(constants2.FlagTendermintStartHeight, constants2.DefaultTendermintStartHeight, fmt.Sprintf("Start checking height on tendermint chain from this height (default %d - starts from where last left)", constants2.DefaultTendermintStartHeight))
	pBridgeCommand.Flags().Int64(constants2.FlagEthereumStartHeight, constants2.DefaultEthereumStartHeight, fmt.Sprintf("Start checking height on ethereum chain from this height (default %d - starts from where last left)", constants2.DefaultEthereumStartHeight))
	pBridgeCommand.Flags().String(constants2.FlagDenom, constants2.DefaultDenom, "denom name")
	pBridgeCommand.Flags().Uint64(constants2.FlagEthGasLimit, constants2.DefaultEthGasLimit, "Gas limit for eth txs")
	pBridgeCommand.Flags().String(constants2.FlagBroadcastMode, constants2.DefaultBroadcastMode, "broadcast mode for tendermint")
	pBridgeCommand.Flags().String(constants2.FlagCASPURL, constants2.DefaultCASPUrl, "casp api url (with http)")
	pBridgeCommand.Flags().String(constants2.FlagCASPVaultID, constants2.DefaultCASPVaultID, "casp vault id")
	pBridgeCommand.Flags().String(constants2.FlagCASPApiToken, constants2.DefaultCASPAPI, "casp api token (in format: Bearer ...)")
	pBridgeCommand.Flags().String(constants2.FlagCASPTMPublicKey, constants2.DefaultCASPTendermintPublicKey, "casp tendermint public key")
	pBridgeCommand.Flags().String(constants2.FlagCASPEthPublicKey, constants2.DefaultCASPEthereumPublicKey, "casp ethereum public key")
	pBridgeCommand.Flags().Int(constants2.FlagCASPSignatureWaitTime, int(constants2.DefaultCASPSignatureWaitTime.Seconds()), "csap siganture wait time")
	pBridgeCommand.Flags().Bool(constants2.FlagCASPConcurrentKey, true, "allows starting multiple sign operations that specify the same key")

	return pBridgeCommand
}

func UpdateConfig(cmd *cobra.Command, pstakeConfig configuration.Config) configuration.Config {
	denom, err := cmd.Flags().GetString(constants2.FlagDenom)
	if err != nil {
		log.Fatalln(err)
	}
	if denom != "" {
		pstakeConfig.Tendermint.PStakeDenom = denom
	}

	ethGasLimit, err := cmd.Flags().GetUint64(constants2.FlagEthGasLimit)
	if err != nil {
		log.Fatalln(err)
	}
	if ethGasLimit != 0 {
		pstakeConfig.Ethereum.EthGasLimit = ethGasLimit
	}

	ports, err := cmd.Flags().GetString(constants2.FlagKafkaPorts)
	if err != nil {
		log.Fatalln(err)
	}
	if ports != "" {
		pstakeConfig.Kafka.Brokers = strings.Split(ports, ",")
	}

	broadcastMode, err := cmd.Flags().GetString(constants2.FlagBroadcastMode)
	if err != nil {
		log.Fatalln(err)
	}
	if broadcastMode == flags.BroadcastBlock || broadcastMode == flags.BroadcastAsync || broadcastMode == flags.BroadcastSync {
		pstakeConfig.Tendermint.BroadcastMode = broadcastMode
	} else {
		log.Fatalln(fmt.Errorf("invalid broadcast mode"))
	}

	caspURL, err := cmd.Flags().GetString(constants2.FlagCASPURL)
	if err != nil {
		log.Fatalln(err)
	}
	if caspURL != "" {
		pstakeConfig.CASP.URL = caspURL
	}

	caspVaultID, err := cmd.Flags().GetString(constants2.FlagCASPVaultID)
	if err != nil {
		log.Fatalln(err)
	}
	if caspVaultID != "" {
		pstakeConfig.CASP.VaultID = caspVaultID
	}

	csapApiToken, err := cmd.Flags().GetString(constants2.FlagCASPApiToken)
	if err != nil {
		log.Fatalln(err)
	}
	if csapApiToken != "" {
		pstakeConfig.CASP.APIToken = csapApiToken
	}

	caspTMPublicKey, err := cmd.Flags().GetString(constants2.FlagCASPTMPublicKey)
	if err != nil {
		log.Fatalln(err)
	}
	if caspTMPublicKey != "" {
		pstakeConfig.CASP.TendermintPublicKey = caspTMPublicKey
	}

	caspEthPublicKey, err := cmd.Flags().GetString(constants2.FlagCASPEthPublicKey)
	if err != nil {
		log.Fatalln(err)
	}
	if caspTMPublicKey != "" {
		pstakeConfig.CASP.EthereumPublicKey = caspEthPublicKey
	}

	caspSignatureWaitTime, err := cmd.Flags().GetInt(constants2.FlagCASPSignatureWaitTime)
	if err != nil {
		log.Fatalln(err)
	}
	if caspSignatureWaitTime >= 0 {
		pstakeConfig.CASP.SignatureWaitTime = time.Duration(caspSignatureWaitTime) * time.Second
	} else {
		log.Fatalln("invalid casp signature wait time")
	}

	caspConcurrentKey, err := cmd.Flags().GetBool(constants2.FlagCASPConcurrentKey)
	if err != nil {
		log.Fatalln(err)
	}
	if caspTMPublicKey != "" {
		pstakeConfig.CASP.AllowConcurrentKeyUsage = caspConcurrentKey
	}

	return pstakeConfig
}
