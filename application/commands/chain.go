package commands

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	ethereum2 "github.com/persistenceOne/persistenceBridge/ethereum"
	"github.com/persistenceOne/persistenceBridge/kafka"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	tendermint2 "github.com/persistenceOne/persistenceBridge/tendermint"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
	"strings"
	"time"
)

func StartCommand(initClientCtx client.Context) *cobra.Command {
	pStakeCommand := &cobra.Command{
		Use:   "start [path_to_chain_json] [mnemonics]",
		Short: "Start persistenceBridge",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			homePath, err := cmd.Flags().GetString(constants2.FlagPStakeHome)
			if err != nil {
				log.Fatalln(err)
			}

			pstakeConfig := configuration.Config{}
			_, err = toml.DecodeFile(filepath.Join(homePath, "config.toml"), &pstakeConfig)
			if err != nil {
				log.Fatalf("Error decoding pstakeConfig file: %v\n", err.Error())
			}

			pstakeConfig = UpdateConfig(cmd, pstakeConfig)

			db, err := application.InitializeDB(homePath+"/db", pstakeConfig.Tendermint.TendermintStartHeight,
				pstakeConfig.Ethereum.EthereumStartHeight)
			if err != nil {
				log.Fatalln(err)
			}
			defer db.Close()

			chain, err := tendermint2.InitializeAndStartChain(args[0], pstakeConfig.Tendermint.RelayerTimeout, homePath, pstakeConfig.CoinType, args[1])
			if err != nil {
				log.Fatalln(err)
			}

			application.SetAppConfiguration(pstakeConfig.PStakeDenom, chain.MustGetAddress(), pstakeConfig.Ethereum.EthAccountPrivateKey,
				pstakeConfig.Ethereum.EthGasLimit)

			ethereumClient, err := ethclient.Dial(pstakeConfig.Ethereum.EthereumEndpoint)
			if err != nil {
				log.Fatalf("Error while dialing to eth orchestrator %s: %s\n", pstakeConfig.Ethereum.EthereumEndpoint, err.Error())
			}

			protoCodec := codec.NewProtoCodec(initClientCtx.InterfaceRegistry)
			kafkaState := utils.NewKafkaState(pstakeConfig.Kafka.Brokers, homePath, pstakeConfig.Kafka.TopicDetail)
			go kafka.KafkaRoutine(kafkaState, pstakeConfig, protoCodec, chain, ethereumClient)
			server.TrapSignal(kafka.KafkaClose(kafkaState))

			log.Println("Starting to listen ethereum....")
			go ethereum2.StartListening(ethereumClient, time.Duration(pstakeConfig.Ethereum.EthereumSleepTime)*time.Millisecond, kafkaState, protoCodec)

			log.Println("Starting to listen tendermint....")
			tendermint2.StartListening(initClientCtx.WithHomeDir(homePath), chain, kafkaState, protoCodec,
				time.Duration(pstakeConfig.Tendermint.TendermintSleepTime)*time.Millisecond)

			return nil
		},
	}
	pStakeCommand.Flags().String(constants2.FlagTimeOut, "", "timeout time for connecting to rpc")
	pStakeCommand.Flags().Uint32(constants2.FlagCoinType, 0, "coin type for wallet")
	pStakeCommand.Flags().String(constants2.FlagPStakeHome, constants2.DefaultPStakeHome, "home for pStake")
	pStakeCommand.Flags().String(constants2.FlagEthereumEndPoint, "", "ethereum orchestrator to connect")
	pStakeCommand.Flags().String("ports", "", "ports kafka brokers are running on, --ports 192.100.10.10:443,192.100.10.11:443")
	pStakeCommand.Flags().Int(constants2.FlagTendermintSleepTime, -1, "sleep time between block checking for tendermint in ms")
	pStakeCommand.Flags().Int(constants2.FlagEthereumSleepTime, -1, "sleep time between block checking for ethereum in ms")
	pStakeCommand.Flags().Int64(constants2.FlagTendermintStartHeight, -1, fmt.Sprintf("Start checking height on tendermint chain from this height (default %d - starts from where last left)", constants2.DefaultTendermintStartHeight))
	pStakeCommand.Flags().Int64(constants2.FlagEthereumStartHeight, -1, fmt.Sprintf("Start checking height on ethereum chain from this height (default %d - starts from where last left)", constants2.DefaultEthereumStartHeight))
	pStakeCommand.Flags().String(constants2.FlagDenom, "", "denom name")
	pStakeCommand.Flags().String(constants2.FlagEthPrivateKey, "", "private keys of ethereum account which does txs.")
	pStakeCommand.Flags().Uint64(constants2.FlagEthGasLimit, 0, "Gas limit for eth txs")
	return pStakeCommand
}

func UpdateConfig(cmd *cobra.Command, pstakeConfig configuration.Config) configuration.Config {
	timeout, err := cmd.Flags().GetString(constants2.FlagTimeOut)
	if err != nil {
		log.Fatalln(err)
	}
	if timeout != "" {
		pstakeConfig.Tendermint.RelayerTimeout = timeout
	}

	coinType, err := cmd.Flags().GetUint32(constants2.FlagCoinType)
	if err != nil {
		log.Fatalln(err)
	}
	if coinType != 0 {
		pstakeConfig.CoinType = coinType
	}

	denom, err := cmd.Flags().GetString(constants2.FlagDenom)
	if err != nil {
		log.Fatalln(err)
	}
	if denom != "" {
		pstakeConfig.PStakeDenom = denom
	}

	tmSleepTime, err := cmd.Flags().GetInt(constants2.FlagTendermintSleepTime)
	if err != nil {
		log.Fatalln(err)
	}
	if tmSleepTime != -1 {
		pstakeConfig.Tendermint.TendermintSleepTime = tmSleepTime
	}
	//tmSleepDuration := time.Duration(tmSleepTime) * time.Millisecond

	ethereumEndPoint, err := cmd.Flags().GetString(constants2.FlagEthereumEndPoint)
	if err != nil {
		log.Fatalln(err)
	}
	if ethereumEndPoint != "" {
		pstakeConfig.Ethereum.EthereumEndpoint = ethereumEndPoint
	}

	ethSleepTime, err := cmd.Flags().GetInt(constants2.FlagEthereumSleepTime)
	if err != nil {
		log.Fatalln(err)
	}
	//ethSleepDuration := time.Duration(ethSleepTime) * time.Millisecond
	if ethSleepTime != -1 {
		pstakeConfig.Ethereum.EthereumSleepTime = ethSleepTime
	}

	tmStart, err := cmd.Flags().GetInt64(constants2.FlagTendermintStartHeight)
	if err != nil {
		log.Fatalln(err)
	}
	if tmStart != -1 {
		pstakeConfig.Tendermint.TendermintStartHeight = tmStart
	}

	ethStart, err := cmd.Flags().GetInt64(constants2.FlagEthereumStartHeight)
	if err != nil {
		log.Fatalln(err)
	}
	if ethStart != -1 {
		pstakeConfig.Ethereum.EthereumStartHeight = ethStart
	}

	ethPrivateKeyStr, err := cmd.Flags().GetString(constants2.FlagEthPrivateKey)
	if err != nil {
		log.Fatalln(err)
	}
	if ethPrivateKeyStr != "" {
		ethPrivateKey, err := crypto.HexToECDSA(ethPrivateKeyStr)
		if err != nil {
			log.Fatal(err)
		}
		pstakeConfig.Ethereum.EthAccountPrivateKey = ethPrivateKey
	}

	ethGasLimit, err := cmd.Flags().GetUint64(constants2.FlagEthGasLimit)
	if err != nil {
		log.Fatalln(err)
	}
	if ethGasLimit != 0 {
		pstakeConfig.Ethereum.EthGasLimit = ethGasLimit
	}

	ports, err := cmd.Flags().GetString("ports")
	if err != nil {
		log.Fatalln(err)
	}
	pstakeConfig.Kafka.Brokers = strings.Split(ports, ",")

	return pstakeConfig
}
