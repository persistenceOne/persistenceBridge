package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/spf13/cobra"
)

func GetCmdWithConfig() *cobra.Command {
	LoadEnv()
	var cmd cobra.Command

	cmd.Flags().String(constants.FlagPBridgeHome, constants.TestHomeDir, "home for pBridge")

	cmd.Flags().String(constants.FlagTendermintNode, os.Getenv("TendermintNode"), "tendermint rpc node url")
	cmd.Flags().String(constants.FlagDenom, constants.DefaultDenom, "denom name")
	cmd.Flags().Uint32(constants.FlagTendermintCoinType, constants.DefaultTendermintCoinType, "tendermint address coin type")
	cmd.Flags().Int64(constants.FlagTMAvgBlockTime, constants.DefaultTendermintAvgBlockTime.Milliseconds(), "avg block of tm chain (in ms)")
	cmd.Flags().String(constants.FlagAccountPrefix, constants.DefaultAccountPrefix, "account prefix on tendermint chains")
	cmd.Flags().String(constants.FlagTendermintChainID, constants.DefaultTendermintChainId, "tendermint rpc node url chains")
	cmd.Flags().String(constants.FlagTMGasPrice, constants.DefaultTendermintGasPrice, "tendermint gas price (should be a float value)")
	cmd.Flags().Float64(constants.FlagTMGasAdjustment, constants.DefaultTendermintGasAdjustment, "tendermint gas adjustment (should be a float value and greater than 1.0)")
	cmd.Flags().String(constants.FlagBroadcastMode, constants.DefaultBroadcastMode, "broadcast mode for tendermint")
	cmd.Flags().Int64(constants.FlagMinimumWrapAmount, constants.DefaultMinimumWrapAmount, "minimum amount in send coin tx to wrap onto eth")

	cmd.Flags().String(constants.FlagEthereumEndPoint, constants.DefaultEthereumEndPoint, "ethereum orchestrator to connect")
	cmd.Flags().Uint64(constants.FlagEthGasLimit, constants.DefaultEthGasLimit, "Gas limit for eth txs")
	cmd.Flags().Int64(constants.FlagEthGasFeeCap, constants.DefaultEthGasFeeCap, "Gas fee cap for eth txs")
	cmd.Flags().String(constants.FlagTokenWrapperAddress, "0x0000000000000000000000000000000000000001", "sc address of token wrapper")
	cmd.Flags().String(constants.FlagLiquidStakingAddress, "0x0000000000000000000000000000000000000001", "sc address of liquid staking")

	cmd.Flags().String(constants.FlagKafkaPorts, constants.DefaultKafkaPorts, "ports kafka brokers are running on, --ports 192.100.10.10:443,192.100.10.11:443")

	cmd.Flags().String(constants.FlagCASPURL, os.Getenv("CASPURL"), "casp api url (with http)")
	cmd.Flags().String(constants.FlagCASPApiToken, os.Getenv("APIToken"), "casp api token")
	cmd.Flags().String(constants.FlagCASPVaultID, os.Getenv("CASPVaultID"), "casp vault id")
	cmd.Flags().String(constants.FlagCASPTMPublicKey, os.Getenv("CASPTMPublicKey"), "casp tendermint public key")
	cmd.Flags().String(constants.FlagCASPEthPublicKey, os.Getenv("CASPEthPublicKey"), "casp ethereum public key")
	cmd.Flags().Int(constants.FlagCASPWaitTime, int(constants.DefaultCASPWaitTime.Seconds()), "casp wait time")
	cmd.Flags().Bool(constants.FlagCASPConcurrentKey, true, "allows starting multiple sign operations that specify the same key")
	cmd.Flags().Uint(constants.FlagCASPMaxAttempts, constants.DefaultCASPMaxAttempts, "max attempts for getting signature")

	cmd.Flags().String(constants.FlagTelegramBotToken, "", "telegram bot token")
	cmd.Flags().Int64(constants.FlagTelegramChatID, 0, "telegram chat id")

	cmd.Flags().Bool(constants.FlagInitSlackBot, true, "slack bot init")

	cmd.Flags().String(constants.FlagRPCEndpoint, constants.DefaultRPCEndpoint, "rpc Endpoint for server")

	return &cmd

}

func LoadEnv() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	err := godotenv.Load(basepath + "/testEnv")
	if err != nil {
		panic(fmt.Errorf("error loading .env file: %v", err))
	}
}
