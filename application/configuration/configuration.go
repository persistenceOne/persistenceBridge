package configuration

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/spf13/cobra"
)

var appConfig *config

func InitConfig() *config {
	c := newConfig()
	appConfig = &c
	return appConfig
}

func GetAppConfig() config {
	return *appConfig
}

func SetPStakeAddress(tmAddress sdk.AccAddress) {
	if !appConfig.seal {
		if strings.Contains(tmAddress.String(), GetAppConfig().Tendermint.AccountPrefix) {
			appConfig.Tendermint.pStakeAddress = tmAddress.String()
		} else {
			panic(fmt.Errorf("pStake wrap address prefix (%s) and Config account prefix (%s) does not match", sdk.GetConfig().GetBech32AccountAddrPrefix(), GetAppConfig().Tendermint.AccountPrefix))
		}
	}
}

func SetConfig(cmd *cobra.Command) *config {
	if appConfig == nil || !appConfig.seal {
		denom, err := cmd.Flags().GetString(constants2.FlagDenom)
		if err != nil {
			log.Fatalln(err)
		}
		if denom != "" {
			appConfig.Tendermint.PStakeDenom = denom
		}

		accountPrefix, err := cmd.Flags().GetString(constants2.FlagAccountPrefix)
		if err != nil {
			log.Fatalln(err)
		}
		if accountPrefix != "" {
			appConfig.Tendermint.AccountPrefix = accountPrefix
		}

		ethereumEndPoint, err := cmd.Flags().GetString(constants2.FlagEthereumEndPoint)
		if err != nil {
			log.Fatalln(err)
		}
		if ethereumEndPoint != "" {
			appConfig.Ethereum.EthereumEndPoint = ethereumEndPoint
		}

		ethGasLimit, err := cmd.Flags().GetUint64(constants2.FlagEthGasLimit)
		if err != nil {
			log.Fatalln(err)
		}
		if ethGasLimit != 0 {
			appConfig.Ethereum.GasLimit = ethGasLimit
		}

		ports, err := cmd.Flags().GetString(constants2.FlagKafkaPorts)
		if err != nil {
			log.Fatalln(err)
		}
		if ports != "" {
			appConfig.Kafka.Brokers = strings.Split(ports, ",")
		}

		broadcastMode, err := cmd.Flags().GetString(constants2.FlagBroadcastMode)
		if err != nil {
			log.Fatalln(err)
		}
		if broadcastMode != "" {
			if broadcastMode == flags.BroadcastBlock || broadcastMode == flags.BroadcastAsync || broadcastMode == flags.BroadcastSync {
				appConfig.Tendermint.BroadcastMode = broadcastMode
			} else {
				log.Fatalln(fmt.Errorf("invalid broadcast mode"))
			}
		}

		caspURL, err := cmd.Flags().GetString(constants2.FlagCASPURL)
		if err != nil {
			log.Fatalln(err)
		}
		if caspURL != "" {
			appConfig.CASP.URL = caspURL
		}

		caspVaultID, err := cmd.Flags().GetString(constants2.FlagCASPVaultID)
		if err != nil {
			log.Fatalln(err)
		}
		if caspVaultID != "" {
			appConfig.CASP.VaultID = caspVaultID
		}

		csapApiToken, err := cmd.Flags().GetString(constants2.FlagCASPApiToken)
		if err != nil {
			log.Fatalln(err)
		}
		if csapApiToken != "" {
			appConfig.CASP.APIToken = csapApiToken
		}

		caspTMPublicKey, err := cmd.Flags().GetString(constants2.FlagCASPTMPublicKey)
		if err != nil {
			log.Fatalln(err)
		}
		if caspTMPublicKey != "" {
			appConfig.CASP.TendermintPublicKey = caspTMPublicKey
		}

		caspEthPublicKey, err := cmd.Flags().GetString(constants2.FlagCASPEthPublicKey)
		if err != nil {
			log.Fatalln(err)
		}
		if caspEthPublicKey != "" {
			appConfig.CASP.EthereumPublicKey = caspEthPublicKey
		}

		caspWaitTime, err := cmd.Flags().GetInt(constants2.FlagCASPWaitTime)
		if err != nil {
			log.Fatalln(err)
		}
		if caspWaitTime > 0 {
			appConfig.CASP.WaitTime = time.Duration(caspWaitTime) * time.Second
		} else if caspWaitTime < 0 {
			log.Fatalln("invalid casp signature wait time")
		}

		caspConcurrentKey, err := cmd.Flags().GetBool(constants2.FlagCASPConcurrentKey)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.CASP.AllowConcurrentKeyUsage = caspConcurrentKey

		caspMaxAttempts, err := cmd.Flags().GetInt(constants2.FlagCASPMaxAttempts)
		if err != nil {
			log.Fatalln(err)
		}
		if caspMaxAttempts > 0 {
			appConfig.CASP.MaxAttempts = caspMaxAttempts
		} else if caspMaxAttempts < 0 {
			log.Fatalln("invalid casp MaxAttempts")
		}

		bridgeRPCEndpoint, err := cmd.Flags().GetString(constants2.FlagRPCEndpoint)
		if err != nil {
			log.Fatalln(err)
		}
		if bridgeRPCEndpoint != "" {
			appConfig.RPCEndpoint = bridgeRPCEndpoint
		}

		minWrapAmt, err := cmd.Flags().GetInt64(constants2.FlagMinimumWrapAmount)
		if err != nil {
			log.Fatalln(err)
		}
		if minWrapAmt >= 0 {
			appConfig.Tendermint.MinimumWrapAmount = minWrapAmt
		}

		telegramBotToken, err := cmd.Flags().GetString(constants2.FlagTelegramBotToken)
		if err != nil {
			log.Fatalln(err)
		}
		if telegramBotToken != "" {
			appConfig.TelegramBot.Token = telegramBotToken
		}

		telegramBotChatID, err := cmd.Flags().GetInt64(constants2.FlagTelegramChatID)
		if err != nil {
			log.Fatalln(err)
		}
		if telegramBotChatID != 0 {
			appConfig.TelegramBot.ChatID = telegramBotChatID
		}
	}

	return appConfig
}

func ValidateAndSeal() {
	if err := appConfig.validate(); err != nil {
		log.Fatalf("configuration validation error: %s", err.Error())
	}
	appConfig.seal = true
}
