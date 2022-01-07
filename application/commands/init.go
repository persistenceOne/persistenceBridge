/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	"bytes"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
)

func InitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "creates config.toml file",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := configuration.SetConfig(cmd)

			var buf bytes.Buffer

			encoder := toml.NewEncoder(&buf)
			if err := encoder.Encode(config); err != nil {
				panic(err)
			}

			homeDir, err := cmd.Flags().GetString(constants.FlagPBridgeHome)
			if err != nil {
				panic(err)
			}

			if err = os.MkdirAll(homeDir, os.ModePerm); err != nil {
				panic(err)
			}

			const writePerm = 0o600

			if err := os.WriteFile(filepath.Join(homeDir, "config.toml"), buf.Bytes(), writePerm); err != nil {
				panic(err)
			}

			log.Println("generated configuration file at ", filepath.Join(homeDir, "config.toml"))

			return nil
		},
	}

	// This will always be used from flag
	cmd.Flags().String(constants.FlagPBridgeHome, constants.DefaultPBridgeHome(), "home for pBridge")

	cmd.Flags().String(constants.FlagEthereumEndPoint, constants.DefaultEthereumEndPoint, "ethereum orchestrator to connect")
	cmd.Flags().String(constants.FlagKafkaPorts, constants.DefaultKafkaPorts, "ports kafka brokers are running on, --ports 192.100.10.10:443,192.100.10.11:443")
	cmd.Flags().String(constants.FlagDenom, constants.DefaultDenom, "denom name")
	cmd.Flags().String(constants.FlagAccountPrefix, constants.DefaultAccountPrefix, "account prefix on tendermint chains")
	cmd.Flags().String(constants.FlagTendermintNode, constants.DefaultTendermintNode, "tendermint rpc node url")
	cmd.Flags().Uint32(constants.FlagTendermintCoinType, constants.DefaultTendermintCoinType, "tendermint address coin type")
	cmd.Flags().String(constants.FlagTendermintChainID, constants.DefaultTendermintChainID, "chain id of tendermint node")
	cmd.Flags().Uint64(constants.FlagEthGasLimit, constants.DefaultEthGasLimit, "Gas limit for eth txs")
	cmd.Flags().String(constants.FlagBroadcastMode, constants.DefaultBroadcastMode, "broadcast mode for tendermint")
	cmd.Flags().String(constants.FlagCASPURL, "", "casp api url (with http)")
	cmd.Flags().String(constants.FlagCASPVaultID, "", "casp vault id")
	cmd.Flags().String(constants.FlagCASPApiToken, "", "casp api token (in format: Bearer ...)")
	cmd.Flags().String(constants.FlagCASPTMPublicKey, "", "casp tendermint public key")
	cmd.Flags().String(constants.FlagCASPEthPublicKey, "", "casp ethereum public key")
	cmd.Flags().Int(constants.FlagCASPSignatureWaitTime, int(constants.DefaultCASPSignatureWaitTime.Seconds()), "casp signature wait time")
	cmd.Flags().Bool(constants.FlagCASPConcurrentKey, true, "allows starting multiple sign operations that specify the same key")
	cmd.Flags().Int(constants.FlagCASPMaxGetSignatureAttempts, constants.DefaultCASPMaxGetSignatureAttempt, "casp max attempts to fetch operation id")
	cmd.Flags().String(constants.FlagRPCEndpoint, constants.DefaultRPCEndpoint, "rpc Endpoint for server")
	cmd.Flags().Int64(constants.FlagMinimumWrapAmount, constants.DefaultMinimumWrapAmount, "minimum amount in send coin tx to wrap onto eth")
	cmd.Flags().String(constants.FlagTelegramBotToken, "", "telegram bot token")
	cmd.Flags().Int64(constants.FlagTelegramChatID, 0, "telegram chat id")

	return cmd
}
