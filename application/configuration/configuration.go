package configuration

import (
	"github.com/BurntSushi/toml"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"path/filepath"
	"strings"
	"time"

	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/spf13/cobra"
)

var appConfig = newConfig()

func GetAppConfig() config {
	return appConfig
}

func InitializeConfigFromFile(homePath string) {
	_, err := toml.DecodeFile(filepath.Join(homePath, "config.toml"), &appConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}
}

func SetConfig(cmd *cobra.Command) config {
	if !appConfig.seal {
		// ---- Tendermint configuration ----
		denom, err := cmd.Flags().GetString(constants2.FlagDenom)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.Tendermint.Denom = denom

		tmAvgTime, err := cmd.Flags().GetInt64(constants2.FlagTMAvgBlockTime)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.Tendermint.AvgBlockTime = time.Duration(tmAvgTime) * time.Millisecond

		accountPrefix, err := cmd.Flags().GetString(constants2.FlagAccountPrefix)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.Tendermint.AccountPrefix = accountPrefix

		tmNode, err := cmd.Flags().GetString(constants2.FlagTendermintNode)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.Tendermint.Node = tmNode

		tmCoinType, err := cmd.Flags().GetUint32(constants2.FlagTendermintCoinType)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.Tendermint.CoinType = tmCoinType

		tmChainID, err := cmd.Flags().GetString(constants2.FlagTendermintChainID)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.Tendermint.ChainID = tmChainID

		tmGasPrice, err := cmd.Flags().GetString(constants2.FlagTMGasPrice)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.Tendermint.GasPrice = tmGasPrice

		tmGasAdjust, err := cmd.Flags().GetFloat64(constants2.FlagTMGasAdjustment)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.Tendermint.GasAdjustment = tmGasAdjust

		broadcastMode, err := cmd.Flags().GetString(constants2.FlagBroadcastMode)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.Tendermint.BroadcastMode = broadcastMode

		minWrapAmt, err := cmd.Flags().GetInt64(constants2.FlagMinimumWrapAmount)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.Tendermint.MinimumWrapAmount = minWrapAmt
		// **** Tendermint configuration ****

		// ---- Ethereum configuration ----
		ethereumEndPoint, err := cmd.Flags().GetString(constants2.FlagEthereumEndPoint)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.Ethereum.EthereumEndPoint = ethereumEndPoint

		ethGasLimit, err := cmd.Flags().GetUint64(constants2.FlagEthGasLimit)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.Ethereum.GasLimit = ethGasLimit

		ethGasFeeCap, err := cmd.Flags().GetInt64(constants2.FlagEthGasFeeCap)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.Ethereum.GasFeeCap = ethGasFeeCap

		tokenWrapperAddress, err := cmd.Flags().GetString(constants2.FlagTokenWrapperAddress)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.Ethereum.TokenWrapperAddress = tokenWrapperAddress

		liquidStakingAddress, err := cmd.Flags().GetString(constants2.FlagLiquidStakingAddress)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.Ethereum.LiquidStakingAddress = liquidStakingAddress
		// **** Ethereum configuration ****

		// ---- Kafka configuration ----
		ports, err := cmd.Flags().GetString(constants2.FlagKafkaPorts)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.Kafka.Brokers = strings.Split(ports, ",")
		// **** Kafka configuration ****

		// ---- CASP configuration ----
		caspURL, err := cmd.Flags().GetString(constants2.FlagCASPURL)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.CASP.URL = caspURL

		caspVaultID, err := cmd.Flags().GetString(constants2.FlagCASPVaultID)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.CASP.VaultID = caspVaultID

		caspTMPublicKey, err := cmd.Flags().GetString(constants2.FlagCASPTMPublicKey)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.CASP.TendermintPublicKey = caspTMPublicKey

		caspEthPublicKey, err := cmd.Flags().GetString(constants2.FlagCASPEthPublicKey)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.CASP.EthereumPublicKey = caspEthPublicKey

		caspWaitTime, err := cmd.Flags().GetInt(constants2.FlagCASPWaitTime)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.CASP.WaitTime = time.Duration(caspWaitTime) * time.Second

		caspConcurrentKey, err := cmd.Flags().GetBool(constants2.FlagCASPConcurrentKey)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.CASP.AllowConcurrentKeyUsage = caspConcurrentKey

		caspMaxAttempts, err := cmd.Flags().GetUint(constants2.FlagCASPMaxAttempts)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.CASP.MaxAttempts = caspMaxAttempts
		// **** CASP configuration ****

		// ---- Telegram configuration ----
		telegramBotToken, err := cmd.Flags().GetString(constants2.FlagTelegramBotToken)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.TelegramBot.Token = telegramBotToken

		telegramBotChatID, err := cmd.Flags().GetInt64(constants2.FlagTelegramChatID)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.TelegramBot.ChatID = telegramBotChatID
		// **** Telegram configuration ****

		bridgeRPCEndpoint, err := cmd.Flags().GetString(constants2.FlagRPCEndpoint)
		if err != nil {
			log.Fatalln(err)
		}
		appConfig.RPCEndpoint = bridgeRPCEndpoint
	}

	return appConfig
}

func ValidateAndSeal() {
	if err := appConfig.validate(); err != nil {
		log.Fatalf("configuration validation error: %s", err.Error())
	}
	appConfig.seal = true
}

func SetCASPAddresses(wrapAddress sdkTypes.AccAddress, bridgeAdminAddress common.Address) {
	if !appConfig.seal {
		setBech32Prefixes()
		setWrapAddress(wrapAddress)
		setBridgeAdminAddress(bridgeAdminAddress)
	}
}

func setBech32Prefixes() {
	if appConfig.Tendermint.AccountPrefix == "" {
		panic("account prefix is empty")
	}
	bech32PrefixAccAddr := appConfig.Tendermint.AccountPrefix
	bech32PrefixAccPub := appConfig.Tendermint.AccountPrefix + sdkTypes.PrefixPublic
	bech32PrefixValAddr := appConfig.Tendermint.AccountPrefix + sdkTypes.PrefixValidator + sdkTypes.PrefixOperator
	bech32PrefixValPub := appConfig.Tendermint.AccountPrefix + sdkTypes.PrefixValidator + sdkTypes.PrefixOperator + sdkTypes.PrefixPublic
	bech32PrefixConsAddr := appConfig.Tendermint.AccountPrefix + sdkTypes.PrefixValidator + sdkTypes.PrefixConsensus
	bech32PrefixConsPub := appConfig.Tendermint.AccountPrefix + sdkTypes.PrefixValidator + sdkTypes.PrefixConsensus + sdkTypes.PrefixPublic

	bech32Configuration := sdkTypes.GetConfig()
	bech32Configuration.SetBech32PrefixForAccount(bech32PrefixAccAddr, bech32PrefixAccPub)
	bech32Configuration.SetBech32PrefixForValidator(bech32PrefixValAddr, bech32PrefixValPub)
	bech32Configuration.SetBech32PrefixForConsensusNode(bech32PrefixConsAddr, bech32PrefixConsPub)
	// Do not seal the config.
}
