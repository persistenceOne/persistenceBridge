package commands

import (
	"bytes"
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

			homeDir, err := cmd.Flags().GetString(constants2.FlagPBridgeHome)
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
	cmd.Flags().String(constants2.FlagDenom, constants2.DefaultDenom, "denom for TM chain")
	cmd.Flags().String(constants2.FlagPBridgeHome, constants2.DefaultPBridgeHome, "home for pBridge")
	cmd.Flags().String(constants2.FlagKafkaPorts, constants2.DefaultKafkaPorts, "ports kafka brokers are running on, --ports 192.100.10.10:443,192.100.10.11:443")
	cmd.Flags().Uint64(constants2.FlagEthGasLimit, constants2.DefaultEthGasLimit, "Gas limit for eth txs")
	cmd.Flags().String(constants2.FlagBroadcastMode, constants2.DefaultBroadcastMode, "broadcast mode for tendermint")
	cmd.Flags().String(constants2.FlagCASPURL, constants2.DefaultCASPUrl, "broadcast mode for tendermint")
	cmd.Flags().String(constants2.FlagCASPVaultID, constants2.DefaultCASPVaultID, "broadcast mode for tendermint")
	cmd.Flags().String(constants2.FlagCASPApiToken, constants2.DefaultCASPAPI, "broadcast mode for tendermint")
	cmd.Flags().String(constants2.FlagCASPTMPublicKey, constants2.DefaultCASPTendermintPublicKey, "broadcast mode for tendermint")
	cmd.Flags().String(constants2.FlagCASPEthPublicKey, constants2.DefaultCASPEthereumPublicKey, "broadcast mode for tendermint")
	cmd.Flags().Int(constants2.FlagCASPSignatureWaitTime, int(constants2.DefaultCASPSignatureWaitTime.Seconds()), "broadcast mode for tendermint")

	return cmd
}
