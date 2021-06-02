package commands

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func InitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "init root command",
		RunE: func(cmd *cobra.Command, args []string) error {

			config := configuration.NewConfig()
			config = UpdateConfig(cmd, config)
			var buf bytes.Buffer
			encoder := toml.NewEncoder(&buf)
			if err := encoder.Encode(config); err != nil {
				panic(err)
			}

			homeDir, err := cmd.Flags().GetString(constants2.FlagPStakeHome)
			if err != nil {
				panic(err)
			}
			if err = os.MkdirAll(homeDir, os.ModePerm); err != nil {
				panic(err)
			}
			if err := ioutil.WriteFile(filepath.Join(homeDir, "config.toml"), buf.Bytes(), 0644); err != nil {
				panic(err)
			}
			log.Println("generated configuration file at ", filepath.Join(homeDir, "config.toml"))

			return nil
		},
	}
	cmd.Flags().String(constants2.FlagTimeOut, constants2.DefaultTimeout, "timeout time for connecting to rpc")
	cmd.Flags().Uint32(constants2.FlagCoinType, constants2.DefaultCoinType, "coin type for wallet")
	cmd.Flags().String(constants2.FlagPStakeHome, constants2.DefaultPStakeHome, "home for pStake")
	cmd.Flags().String(constants2.FlagEthereumEndPoint, constants2.DefaultEthereumEndPoint, "ethereum orchestrator to connect")
	cmd.Flags().String("ports", "localhost:9092", "ports kafka brokers are running on, --ports 192.100.10.10:443,192.100.10.11:443")
	cmd.Flags().Int(constants2.FlagTendermintSleepTime, constants2.DefaultTendermintSleepTime, "sleep time between block checking for tendermint in ms")
	cmd.Flags().Int(constants2.FlagEthereumSleepTime, constants2.DefaultEthereumSleepTime, "sleep time between block checking for ethereum in ms")
	cmd.Flags().Int64(constants2.FlagTendermintStartHeight, constants2.DefaultTendermintStartHeight, fmt.Sprintf("Start checking height on tendermint chain from this height (default %d - starts from where last left)", constants2.DefaultTendermintStartHeight))
	cmd.Flags().Int64(constants2.FlagEthereumStartHeight, constants2.DefaultEthereumStartHeight, fmt.Sprintf("Start checking height on ethereum chain from this height (default %d - starts from where last left)", constants2.DefaultEthereumStartHeight))
	cmd.Flags().String(constants2.FlagDenom, constants2.DefaultDenom, "denom name")
	cmd.Flags().String(constants2.FlagEthPrivateKey, "", "private keys of ethereum account which does txs.")
	cmd.Flags().Uint64(constants2.FlagEthGasLimit, constants2.DefaultEthGasLimit, "Gas limit for eth txs")

	return cmd
}
