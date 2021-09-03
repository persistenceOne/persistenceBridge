package configuration

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetAppConfig(t *testing.T) {
	t.Logf("Testing GetAppConfig()")
	var Config = GetAppConfig()
	require.Equal(t, appConfig , Config, "The two configurations should be the same")
}

func TestInitConfig(t *testing.T) {
	t.Logf("Testing InitConfig()")
	newAppConfigAddress := InitConfig()
	newAppConfig := *newAppConfigAddress
	appConfig := newConfig()

	require.Equal(t, appConfig, newAppConfig)

	require.Equal(t, appConfig.Kafka, newKafkaConfig(), "Kafka config not initalized properly")

	require.Equal(t, appConfig.Ethereum, newEthereumConfig(), "Ethereum config not initalized properly")

	require.Equal(t, appConfig.TelegramBot, newTelegramBot(), "TelegramBot config not initalized properly")

	require.Equal(t, appConfig.seal, false, "seal not initalized properly")

	require.Equal(t, appConfig.RPCEndpoint, constants2.DefaultRPCEndpoint, "RPCEndpoints not initalized properly")

	require.Equal(t, appConfig.Tendermint.PStakeDenom, newTendermintConfig().PStakeDenom, "PStakeDenom not initalized properly")

	require.Equal(t, appConfig.Tendermint.MinimumWrapAmount, newTendermintConfig().MinimumWrapAmount, "MinimumWrapAmount not initalized properly")

	require.Equal(t, appConfig.Tendermint.BroadcastMode, newTendermintConfig().BroadcastMode, "BroadcastMode not initalized properly")

	require.Equal(t, appConfig.CASP.SignatureWaitTime, constants2.DefaultCASPSignatureWaitTime, "BroadcastMode not initalized properly")

	require.Equal(t, appConfig.CASP.AllowConcurrentKeyUsage, true, "BroadcastMode not initalized properly")

}

func TestSetConfig(t *testing.T) {
	appConfig = InitConfig()
	var cmd cobra.Command

	cmd.Flags().String(constants2.FlagPBridgeHome, constants2.DefaultPBridgeHome, "home for pBridge")
	cmd.Flags().String(constants2.FlagEthereumEndPoint, constants2.DefaultEthereumEndPoint, "ethereum orchestrator to connect")
	cmd.Flags().String(constants2.FlagKafkaPorts, constants2.DefaultKafkaPorts, "ports kafka brokers are running on, --ports 192.100.10.10:443,192.100.10.11:443")
	cmd.Flags().String(constants2.FlagDenom, constants2.DefaultDenom, "denom name")
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
	cmd.Flags().String(constants2.FlagTelegramBotToken, "xyz", "telegram bot token")
	cmd.Flags().Int64(constants2.FlagTelegramChatID, 0, "telegram chat id")
	cmd.Flags().Int(constants2.FlagCASPMaxGetSignatureAttempts, constants2.DefaultCASPMaxGetSignatureAttempt, "max attempts for getting signature")

	appConfig := SetConfig(&cmd)

	denom, _ := cmd.Flags().GetString(constants2.FlagDenom)
	require.Equal(t, appConfig.Tendermint.PStakeDenom, denom, "FlagDenom not set")

	ethereumEndPoint, _ := cmd.Flags().GetString(constants2.FlagEthereumEndPoint)
	require.Equal(t, appConfig.Ethereum.EthereumEndPoint, ethereumEndPoint, "EthereumEndPoint not set")

	ethGasLimit, _ := cmd.Flags().GetUint64(constants2.FlagEthGasLimit)
	require.Equal(t, appConfig.Ethereum.GasLimit, ethGasLimit, "GasLimit not set")

	ports, _ := cmd.Flags().GetString(constants2.FlagKafkaPorts)
	require.Equal(t, appConfig.Kafka.Brokers[0], ports, "Brokers not set")

	broadcastMode, _ := cmd.Flags().GetString(constants2.FlagBroadcastMode)
	require.Equal(t, appConfig.Tendermint.BroadcastMode, broadcastMode, "BroadcastMode not set")

	caspURL, _ := cmd.Flags().GetString(constants2.FlagCASPURL)
	require.Equal(t, appConfig.CASP.URL, caspURL, "CASP URL not set")

	caspVaultID, _ := cmd.Flags().GetString(constants2.FlagCASPVaultID)
	require.Equal(t, appConfig.CASP.VaultID, caspVaultID, "CASP VaultID not set")

	caspApiToken, _ := cmd.Flags().GetString(constants2.FlagCASPApiToken)
	require.Equal(t, appConfig.CASP.APIToken, caspApiToken, "CASP APItoken not set")

	caspTMPublicKey, _ := cmd.Flags().GetString(constants2.FlagCASPTMPublicKey)
	require.Equal(t, appConfig.CASP.TendermintPublicKey, caspTMPublicKey, "TendermintPublicKey not set")

	caspEthPublicKey, _ := cmd.Flags().GetString(constants2.FlagCASPEthPublicKey)
	require.Equal(t, appConfig.CASP.EthereumPublicKey, caspEthPublicKey, "EthereumPublicKey not set")

	caspSignatureWaitTime, _ := cmd.Flags().GetInt(constants2.FlagCASPSignatureWaitTime)
	require.Equal(t, int(appConfig.CASP.SignatureWaitTime/time.Second), caspSignatureWaitTime, "CASP signature wait time not set")

	caspConcurrentKey, _ := cmd.Flags().GetBool(constants2.FlagCASPConcurrentKey)
	require.Equal(t, appConfig.CASP.AllowConcurrentKeyUsage, caspConcurrentKey, "AllowConcurrentKeyUsage not set")

	bridgeRPCEndpoint, _ := cmd.Flags().GetString(constants2.FlagRPCEndpoint)
	require.Equal(t, appConfig.RPCEndpoint, bridgeRPCEndpoint, "RPCEndpoint not set")

	minWrapAmt, _ := cmd.Flags().GetInt64(constants2.FlagMinimumWrapAmount)
	require.Equal(t, appConfig.Tendermint.MinimumWrapAmount, minWrapAmt, "MinimumWrapAmount not set")

	telegramBotToken, _ := cmd.Flags().GetString(constants2.FlagTelegramBotToken)
	require.Equal(t, appConfig.TelegramBot.Token, telegramBotToken, "Telegram Token not set")

	telegramBotChatID, _ := cmd.Flags().GetInt64(constants2.FlagTelegramChatID)
	require.Equal(t, appConfig.TelegramBot.ChatID, telegramBotChatID, "ChatId nots set")

	caspSignatureAttempts, _ := cmd.Flags().GetInt(constants2.FlagCASPMaxGetSignatureAttempts)
	require.Equal(t, appConfig.CASP.MaxGetSignatureAttempts, caspSignatureAttempts, "CASP signature wait time not set")

}

func TestSetPStakeAddress(t *testing.T) {
	InitConfig()
	appConfig := test.GetCmdWithConfig()
	config := SetConfig(&appConfig)
	opaddress, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	SetPStakeAddress(opaddress)
	require.Equal(t, config.Tendermint.pStakeAddress, opaddress.String(), "PStakeAddress not set")
}

func TestValidateAndSeal(t *testing.T) {
	InitConfig()
	appConfig := test.GetCmdWithConfig()
	config := SetConfig(&appConfig)
	opaddress, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	SetPStakeAddress(opaddress)
	ValidateAndSeal()
	require.Equal(t, config.seal, true, "appConfig did not get validated")
}

func TestGetPStakeAddress(t *testing.T) {
	InitConfig()
	appConfig := test.GetCmdWithConfig()
	config := SetConfig(&appConfig)
	opaddress, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	SetPStakeAddress(opaddress)
	PstakeAddress := config.Tendermint.GetPStakeAddress()
	require.Equal(t, PstakeAddress, opaddress.String(), "pStakeAddress not set")
}
