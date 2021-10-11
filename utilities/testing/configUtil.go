package testing

import (
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/spf13/cobra"
	"os"
)

func GetCmdWithConfig() *cobra.Command {
	var cmd cobra.Command

	cmd.Flags().String(constants2.FlagPBridgeHome, constants2.TestHomeDir, "home for pBridge")
	cmd.Flags().String(constants2.FlagPBridgeHome, constants2.DefaultPBridgeHome, "home for pBridge")
	cmd.Flags().String(constants2.FlagEthereumEndPoint, constants2.DefaultEthereumEndPoint, "ethereum orchestrator to connect")
	cmd.Flags().String(constants2.FlagKafkaPorts, constants2.DefaultKafkaPorts, "ports kafka brokers are running on, --ports 192.100.10.10:443,192.100.10.11:443")
	cmd.Flags().String(constants2.FlagDenom, constants2.DefaultDenom, "denom name")
	cmd.Flags().String(constants2.FlagAccountPrefix, constants2.DefaultAccountPrefix, "account prefix on tendermint chains")
	cmd.Flags().String(constants2.FlagTendermintNode, constants2.DefaultTendermintNode, "tendermint rpc node url")
	cmd.Flags().String(constants2.FlagTendermintChainID, constants2.DefaultTendermintChainId, "tendermint rpc node url chains")
	cmd.Flags().String(constants2.FlagTMGasPrice, constants2.DefaultTendermintGasPrice, "tendermint gas price (should be a float value)")
	cmd.Flags().Float64(constants2.FlagTMGasAdjustment, constants2.DefaultTendermintGasAdjustment, "tendermint gas adjustment (should be a float value and greater than 1.0)")
	cmd.Flags().Uint64(constants2.FlagEthGasLimit, constants2.DefaultEthGasLimit, "Gas limit for eth txs")
	cmd.Flags().Int64(constants2.FlagEthGasFeeCap, constants2.DefaultEthGasFeeCap, "Gas fee cap for eth txs")
	cmd.Flags().String(constants2.FlagTokenWrapperAddress, constants2.DefaultEthZeroAddress, "sc address of token wrapper")
	cmd.Flags().String(constants2.FlagLiquidStakingAddress, constants2.DefaultEthZeroAddress, "sc address of liquid staking")
	cmd.Flags().String(constants2.FlagBroadcastMode, constants2.DefaultBroadcastMode, "broadcast mode for tendermint")
	cmd.Flags().String(constants2.FlagCASPURL, os.Getenv("CASPURL"), "casp api url (with http)")
	cmd.Flags().String(constants2.FlagCASPVaultID, os.Getenv("CASPVaultID"), "casp vault id")
	cmd.Flags().String(constants2.FlagCASPTMPublicKey, os.Getenv("CASPTMPublicKey"), "casp tendermint public key")
	cmd.Flags().String(constants2.FlagCASPEthPublicKey, os.Getenv("CASPEthPublicKey"), "casp ethereum public key")
	cmd.Flags().Int(constants2.FlagCASPWaitTime, int(constants2.DefaultCASPWaitTime.Seconds()), "casp wait time")
	cmd.Flags().Bool(constants2.FlagCASPConcurrentKey, true, "allows starting multiple sign operations that specify the same key")
	cmd.Flags().String(constants2.FlagRPCEndpoint, constants2.DefaultRPCEndpoint, "rpc Endpoint for server")
	cmd.Flags().Int64(constants2.FlagMinimumWrapAmount, constants2.DefaultMinimumWrapAmount, "minimum amount in send coin tx to wrap onto eth")
	cmd.Flags().String(constants2.FlagTelegramBotToken, "", "telegram bot token")
	cmd.Flags().Int64(constants2.FlagTelegramChatID, 0, "telegram chat id")
	cmd.Flags().Int(constants2.FlagCASPMaxAttempts, constants2.DefaultCASPMaxAttempts, "max attempts for getting signature")

	return &cmd

}
