/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package testing

import (
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/spf13/cobra"
)

func GetCmdWithConfig() *cobra.Command {
	var cmd cobra.Command

	cmd.Flags().String(constants2.FlagPBridgeHome, constants2.DefaultPBridgeHome, "home for pBridge")
	cmd.Flags().String(constants2.FlagEthereumEndPoint, constants2.DefaultEthereumEndPoint, "ethereum orchestrator to connect")
	cmd.Flags().String(constants2.FlagKafkaPorts, constants2.DefaultKafkaPorts, "ports kafka brokers are running on, --ports 192.100.10.10:443,192.100.10.11:443")
	cmd.Flags().String(constants2.FlagDenom, constants2.DefaultDenom, "denom name")
	cmd.Flags().String(constants2.FlagAccountPrefix, constants2.DefaultAccountPrefix, "account prefix on tendermint chains")
	cmd.Flags().String(constants2.FlagTendermintNode, constants2.DefaultTendermintNode, "tendermint rpc node url")
	cmd.Flags().String(constants2.FlagTendermintChainID, constants2.DefaultTendermintChainId, "tendermint rpc node url chains")
	cmd.Flags().Uint64(constants2.FlagEthGasLimit, constants2.DefaultEthGasLimit, "Gas limit for eth txs")
	cmd.Flags().String(constants2.FlagBroadcastMode, constants2.DefaultBroadcastMode, "broadcast mode for tendermint")
	cmd.Flags().String(constants2.FlagCASPURL, "https://65.2.149.241:443", "casp api url (with http)")
	cmd.Flags().String(constants2.FlagCASPVaultID, "4ec017bf-4af8-41b3-9527-a466e05971cb", "casp vault id")
	cmd.Flags().String(constants2.FlagCASPApiToken, "Bearer cHVuZWV0TmV3QXBpa2V5MTI6OWM1NDBhMzAtNTQ5NC00ZDdhLTljODktODA3MDZiNWNhYzQ1", "casp api token (in format: Bearer ...)")
	cmd.Flags().String(constants2.FlagCASPTMPublicKey, "3056301006072A8648CE3D020106052B8104000A0342000413109ECEADCBF6122EF44184B207F8C6820E509497792DDFB166BC090A0FB4447CFFCE16BAAF9EC7F57D14C02641B3A6A698614D973ED744E725A85E62535DA4", "casp tendermint public key")
	cmd.Flags().String(constants2.FlagCASPEthPublicKey, "3056301006072A8648CE3D020106052B8104000A034200049D8BB9DC3E37511273286F60C989BFFC3E28909F426AF7D4A7899FACC4E3DB00413E2DA7A8CF33F367D8C4D8FC2BFA791DD4389CC1E75154CD38429FD9525946", "casp ethereum public key")
	cmd.Flags().Int(constants2.FlagCASPSignatureWaitTime, int(constants2.DefaultCASPSignatureWaitTime.Seconds()), "casp signature wait time")
	cmd.Flags().Bool(constants2.FlagCASPConcurrentKey, true, "allows starting multiple sign operations that specify the same key")
	cmd.Flags().String(constants2.FlagRPCEndpoint, constants2.DefaultRPCEndpoint, "rpc Endpoint for server")
	cmd.Flags().Int64(constants2.FlagMinimumWrapAmount, constants2.DefaultMinimumWrapAmount, "minimum amount in send coin tx to wrap onto eth")
	cmd.Flags().String(constants2.FlagTelegramBotToken, "", "telegram bot token")
	cmd.Flags().Int64(constants2.FlagTelegramChatID, 0, "telegram chat id")
	cmd.Flags().Int(constants2.FlagCASPMaxGetSignatureAttempts, constants2.DefaultCASPMaxGetSignatureAttempt, "max attempts for getting signature")

	return &cmd

}
