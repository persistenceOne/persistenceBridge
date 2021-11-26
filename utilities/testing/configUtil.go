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
	cmd.Flags().String(constants2.FlagCASPURL, "https://casp.test.persistence.one/", "casp api url (with http)")
	cmd.Flags().String(constants2.FlagCASPVaultID, "b9b23cfd-2c53-45c4-a949-d07f7348b910", "casp vault id")
	cmd.Flags().String(constants2.FlagCASPApiToken, "Bearer Z29rdWw6MjBkZTM3OTktYjI0NC00YzIxLWIzMDEtYWRlZDk2MWI4Nzc2", "casp api token (in format: Bearer ...)")
	cmd.Flags().String(constants2.FlagCASPTMPublicKey, "3056301006072A8648CE3D020106052B8104000A034200044AD9710EF4AEA73C200DC6F08B792CDBCE5EE89153FA986C22B4A527BADF0EBDACBDD00A67C0F3A331BFC35D65CBBEAA8509EFFFFE1EF645A5C4EB15BF648D37", "casp tendermint public key")
	cmd.Flags().String(constants2.FlagCASPEthPublicKey, "3056301006072A8648CE3D020106052B8104000A03420004B579EE6C480AB79820975E17211181FCA3E06B780D898D3CB181CD4B02984479DB4E5F1F0985E976251D5DE8B819E6BE86C59F8ED6CCA711B1F7772F09E5B97B", "casp ethereum public key")
	cmd.Flags().Int(constants2.FlagCASPSignatureWaitTime, int(constants2.DefaultCASPSignatureWaitTime.Seconds()), "casp signature wait time")
	cmd.Flags().Bool(constants2.FlagCASPConcurrentKey, true, "allows starting multiple sign operations that specify the same key")
	cmd.Flags().String(constants2.FlagRPCEndpoint, constants2.DefaultRPCEndpoint, "rpc Endpoint for server")
	cmd.Flags().Int64(constants2.FlagMinimumWrapAmount, constants2.DefaultMinimumWrapAmount, "minimum amount in send coin tx to wrap onto eth")
	cmd.Flags().String(constants2.FlagTelegramBotToken, "", "telegram bot token")
	cmd.Flags().Int64(constants2.FlagTelegramChatID, 0, "telegram chat id")
	cmd.Flags().Int(constants2.FlagCASPMaxGetSignatureAttempts, constants2.DefaultCASPMaxGetSignatureAttempt, "max attempts for getting signature")

	return &cmd

}
