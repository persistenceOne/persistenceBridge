package commands

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	ethereum2 "github.com/persistenceOne/persistenceBridge/ethereum"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	tendermint2 "github.com/persistenceOne/persistenceBridge/tendermint"
	"github.com/spf13/cobra"
	"log"
	"strings"
	"time"
)

func StartCommand(initClientCtx client.Context) *cobra.Command {
	pStakeCommand := &cobra.Command{
		Use:   "start [path_to_chain_json] [mnemonics]",
		Short: "Start persistenceBridge",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			timeout, err := cmd.Flags().GetString(constants2.FlagTimeOut)
			if err != nil {
				log.Fatalln(err)
			}

			coinType, err := cmd.Flags().GetUint32(constants2.FlagCoinType)
			if err != nil {
				log.Fatalln(err)
			}

			denom, err := cmd.Flags().GetString(constants2.FlagDenom)
			if err != nil {
				log.Fatalln(err)
			}

			homePath, err := cmd.Flags().GetString(constants2.FlagPStakeHome)
			if err != nil {
				log.Fatalln(err)
			}

			tmSleepTime, err := cmd.Flags().GetInt(constants2.FlagTendermintSleepTime)
			if err != nil {
				log.Fatalln(err)
			}
			tmSleepDuration := time.Duration(tmSleepTime) * time.Millisecond

			ethereumEndPoint, err := cmd.Flags().GetString(constants2.FlagEthereumEndPoint)
			if err != nil {
				log.Fatalln(err)
			}

			ethSleepTime, err := cmd.Flags().GetInt(constants2.FlagEthereumSleepTime)
			if err != nil {
				log.Fatalln(err)
			}
			ethSleepDuration := time.Duration(ethSleepTime) * time.Millisecond

			ports, err := cmd.Flags().GetString("ports")
			if err != nil {
				log.Fatalln(err)
			}

			kafkaHome, err := cmd.Flags().GetString(utils.FlagKafkaHome)
			if err != nil {
				log.Fatalln(err)
			}

			tmStart, err := cmd.Flags().GetInt64(constants2.FlagTendermintStartHeight)
			if err != nil {
				log.Fatalln(err)
			}

			ethStart, err := cmd.Flags().GetInt64(constants2.FlagEthereumStartHeight)
			if err != nil {
				log.Fatalln(err)
			}

			ethPrivateKeyStr, err := cmd.Flags().GetString(constants2.FlagEthPrivateKey)
			if err != nil {
				log.Fatalln(err)
			}
			ethPrivateKey, err := crypto.HexToECDSA(ethPrivateKeyStr)
			if err != nil {
				log.Fatal(err)
			}

			ethGasLimit, err := cmd.Flags().GetUint64(constants2.FlagEthGasLimit)
			if err != nil {
				log.Fatalln(err)
			}

			db, err := application.InitializeDB(homePath+"/db", tmStart, ethStart)
			if err != nil {
				log.Fatalln(err)
			}
			defer db.Close()

			chain, err := tendermint2.InitializeAndStartChain(args[0], timeout, homePath, coinType, args[1])
			if err != nil {
				log.Fatalln(err)
			}

			application.SetAppConfiguration(denom, chain.MustGetAddress(), ethPrivateKey, ethGasLimit)

			ethereumClient, err := ethclient.Dial(ethereumEndPoint)
			if err != nil {
				log.Fatalf("Error while dialing to eth orchestrator %s: %s\n", ethereumEndPoint, err.Error())
			}

			protoCodec := codec.NewProtoCodec(initClientCtx.InterfaceRegistry)
			portsList := strings.Split(ports, ",")
			kafkaState := utils.NewKafkaState(portsList, kafkaHome)
			go kafkaRoutine(kafkaState, protoCodec, chain, ethereumClient)
			server.TrapSignal(kafkaClose(kafkaState))

			log.Println("Starting to listen ethereum....")
			go ethereum2.StartListening(ethereumClient, ethSleepDuration, kafkaState, protoCodec)

			log.Println("Starting to listen tendermint....")
			tendermint2.StartListening(initClientCtx.WithHomeDir(homePath), chain, kafkaState, protoCodec, tmSleepDuration)

			return nil
		},
	}
	pStakeCommand.Flags().String(constants2.FlagTimeOut, constants2.DefaultTimeout, "timeout time for connecting to rpc")
	pStakeCommand.Flags().Uint32(constants2.FlagCoinType, constants2.DefaultCoinType, "coin type for wallet")
	pStakeCommand.Flags().String(constants2.FlagPStakeHome, constants2.DefaultPStakeHome, "home for pStake")
	pStakeCommand.Flags().String(constants2.FlagEthereumEndPoint, constants2.DefaultEthereumEndPoint, "ethereum orchestrator to connect")
	pStakeCommand.Flags().String("ports", "localhost:9092", "ports kafka brokers are running on, --ports 192.100.10.10:443,192.100.10.11:443")
	pStakeCommand.Flags().String(utils.FlagKafkaHome, utils.DefaultKafkaHome, "The kafka configuration file directory")
	pStakeCommand.Flags().Int(constants2.FlagTendermintSleepTime, constants2.DefaultTendermintSleepTime, "sleep time between block checking for tendermint in ms")
	pStakeCommand.Flags().Int(constants2.FlagEthereumSleepTime, constants2.DefaultEthereumSleepTime, "sleep time between block checking for ethereum in ms")
	pStakeCommand.Flags().Int64(constants2.FlagTendermintStartHeight, constants2.DefaultTendermintStartHeight, fmt.Sprintf("Start checking height on tendermint chain from this height (default %d - starts from where last left)", constants2.DefaultTendermintStartHeight))
	pStakeCommand.Flags().Int64(constants2.FlagEthereumStartHeight, constants2.DefaultEthereumStartHeight, fmt.Sprintf("Start checking height on ethereum chain from this height (default %d - starts from where last left)", constants2.DefaultEthereumStartHeight))
	pStakeCommand.Flags().String(constants2.FlagDenom, constants2.DefaultDenom, "denom name")
	pStakeCommand.Flags().String(constants2.FlagEthPrivateKey, "", "private keys of ethereum account which does txs.")
	pStakeCommand.Flags().Uint64(constants2.FlagEthGasLimit, constants2.DefaultEthGasLimit, "Gas limit for eth txs")
	return pStakeCommand
}
