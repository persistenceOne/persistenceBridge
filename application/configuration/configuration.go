/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package configuration

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/persistenceOne/persistenceBridge/application/constants"
)

// fixme: we cannot use such things as global state if we'd like to have parallel tests
var appConfig = newConfig()

func GetAppConfig() config {
	return appConfig.DeepCopy()
}

func SetPStakeAddress(tmAddress sdk.AccAddress) {
	if !appConfig.seal {
		if strings.Contains(tmAddress.String(), GetAppConfig().Tendermint.AccountPrefix) {
			appConfig.Tendermint.pStakeAddress = tmAddress.String()
		} else {
			panic(fmt.Errorf("%w: address prefix (%s), config account prefix (%s)",
				ErrIncorrectAccountPrefix, sdk.GetConfig().GetBech32AccountAddrPrefix(), GetAppConfig().Tendermint.AccountPrefix))
		}
	}
}

func SetConfig(cmd *cobra.Command) config {
	if appConfig.isFilled || appConfig.seal {
		return appConfig
	}

	denom, err := cmd.Flags().GetString(constants.FlagDenom)
	if err != nil {
		log.Fatalln(err)
	}

	if denom != "" {
		appConfig.Tendermint.PStakeDenom = denom
	}

	accountPrefix, err := cmd.Flags().GetString(constants.FlagAccountPrefix)
	if err != nil {
		log.Fatalln(err)
	}

	if accountPrefix != "" {
		appConfig.Tendermint.AccountPrefix = accountPrefix
	}

	ethereumEndPoint, err := cmd.Flags().GetString(constants.FlagEthereumEndPoint)
	if err != nil {
		log.Fatalln(err)
	}

	if ethereumEndPoint != "" {
		appConfig.Ethereum.EthereumEndPoint = ethereumEndPoint
	}

	ethGasLimit, err := cmd.Flags().GetUint64(constants.FlagEthGasLimit)
	if err != nil {
		log.Fatalln(err)
	}

	if ethGasLimit != 0 {
		appConfig.Ethereum.GasLimit = ethGasLimit
	}

	ports, err := cmd.Flags().GetString(constants.FlagKafkaPorts)
	if err != nil {
		log.Fatalln(err)
	}

	if ports != "" {
		appConfig.Kafka.Brokers = strings.Split(ports, ",")
	}

	broadcastMode, err := cmd.Flags().GetString(constants.FlagBroadcastMode)
	if err != nil {
		log.Fatalln(err)
	}

	if broadcastMode != "" {
		if broadcastMode == flags.BroadcastBlock || broadcastMode == flags.BroadcastAsync || broadcastMode == flags.BroadcastSync {
			appConfig.Tendermint.BroadcastMode = broadcastMode
		} else {
			log.Fatalln(ErrInvalidBroadcastMode)
		}
	}

	caspURL, err := cmd.Flags().GetString(constants.FlagCASPURL)
	if err != nil {
		log.Fatalln(err)
	}

	if caspURL != "" {
		appConfig.CASP.URL = caspURL
	}

	caspVaultID, err := cmd.Flags().GetString(constants.FlagCASPVaultID)
	if err != nil {
		log.Fatalln(err)
	}

	if caspVaultID != "" {
		appConfig.CASP.VaultID = caspVaultID
	}

	csapAPIToken, err := cmd.Flags().GetString(constants.FlagCASPApiToken)
	if err != nil {
		log.Fatalln(err)
	}

	if csapAPIToken != "" {
		appConfig.CASP.APIToken = csapAPIToken
	}

	caspTMPublicKey, err := cmd.Flags().GetString(constants.FlagCASPTMPublicKey)
	if err != nil {
		log.Fatalln(err)
	}

	if caspTMPublicKey != "" {
		appConfig.CASP.TendermintPublicKey = caspTMPublicKey
	}

	caspEthPublicKey, err := cmd.Flags().GetString(constants.FlagCASPEthPublicKey)
	if err != nil {
		log.Fatalln(err)
	}

	if caspTMPublicKey != "" {
		appConfig.CASP.EthereumPublicKey = caspEthPublicKey
	}

	caspSignatureWaitTime, err := cmd.Flags().GetInt(constants.FlagCASPSignatureWaitTime)
	if err != nil {
		log.Fatalln(err)
	}

	if caspSignatureWaitTime >= 0 {
		appConfig.CASP.SignatureWaitTime = time.Duration(caspSignatureWaitTime) * time.Second
	} else if appConfig.CASP.SignatureWaitTime < 0 {
		log.Fatalln("invalid casp signature wait time")
	}

	caspConcurrentKey, err := cmd.Flags().GetBool(constants.FlagCASPConcurrentKey)
	if err != nil {
		log.Fatalln(err)
	}

	appConfig.CASP.AllowConcurrentKeyUsage = caspConcurrentKey

	caspMaxGetSignatureAttempts, err := cmd.Flags().GetInt(constants.FlagCASPMaxGetSignatureAttempts)
	if err != nil {
		log.Fatalln(err)
	}

	if caspMaxGetSignatureAttempts > 0 {
		appConfig.CASP.MaxGetSignatureAttempts = caspMaxGetSignatureAttempts
	} else if appConfig.CASP.SignatureWaitTime < 0 {
		log.Fatalln("invalid casp MaxGetSignatureAttempts")
	}

	bridgeRPCEndpoint, err := cmd.Flags().GetString(constants.FlagRPCEndpoint)
	if err != nil {
		log.Fatalln(err)
	}

	if bridgeRPCEndpoint != "" {
		appConfig.RPCEndpoint = bridgeRPCEndpoint
	}

	minWrapAmt, err := cmd.Flags().GetInt64(constants.FlagMinimumWrapAmount)
	if err != nil {
		log.Fatalln(err)
	}

	if minWrapAmt >= 0 {
		appConfig.Tendermint.MinimumWrapAmount = minWrapAmt
	}

	telegramBotToken, err := cmd.Flags().GetString(constants.FlagTelegramBotToken)
	if err != nil {
		log.Fatalln(err)
	}

	if telegramBotToken != "" {
		appConfig.TelegramBot.Token = telegramBotToken
	}

	telegramBotChatID, err := cmd.Flags().GetInt64(constants.FlagTelegramChatID)
	if err != nil {
		log.Fatalln(err)
	}

	if telegramBotChatID != 0 {
		appConfig.TelegramBot.ChatID = telegramBotChatID
	}

	appConfig.isFilled = true

	return appConfig
}

func ValidateAndSeal() {
	if err := appConfig.validate(); err != nil {
		log.Fatalf("configuration validation error: %s", err.Error())
	}

	appConfig.seal = true
}
