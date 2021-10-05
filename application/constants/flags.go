package constants

import (
	"github.com/cosmos/cosmos-sdk/client/flags"
	"os"
	"time"
)

const (
	FlagTimeOut               = "timeout"
	FlagPBridgeHome           = "pBridgeHome"
	FlagEthereumEndPoint      = "ethEndPoint"
	FlagTendermintSleepTime   = "tmSleep"
	FlagEthereumSleepTime     = "ethSleep"
	FlagTendermintStartHeight = "tmStart"
	FlagEthereumStartHeight   = "ethStart"
	FlagDenom                 = "denom"
	FlagTendermintCoinType    = "tmCoinType"
	FlagEthGasLimit           = "ethGasLimit"
	FlagKafkaPorts            = "kafkaPorts"
	FlagBroadcastMode         = "tmBroadcastMode"
	FlagCASPURL               = "caspURL"
	FlagCASPVaultID           = "caspVaultID"
	FlagCASPTMPublicKey       = "caspTMPublicKey"
	FlagCASPEthPublicKey      = "caspEthPublicKey"
	FlagCASPWaitTime          = "caspWaitTime"
	FlagCASPConcurrentKey     = "caspConcurrentKeyUsage"
	FlagCASPMaxAttempts       = "caspMaxAttempts"
	FlagRPCEndpoint           = "rpc-endpoint"
	FlagMinimumWrapAmount     = "minWrapAmt"
	FlagTelegramBotToken      = "botToken"
	FlagTelegramChatID        = "botChatID"
	FlagShowDebugLog          = "debug"
	FlagAccountPrefix         = "accountPrefix"
	FlagTendermintNode        = "tmNode"
	FlagTendermintChainID     = "chainID"
)

var (
	DefaultTimeout                 = "10s"
	DefaultEthereumEndPoint        = "wss://ropsten.infura.io/ws/v3/b21966541db246d398fb31402eec2c14"
	DefaultTendermintSleepTime     = 3000     //ms
	DefaultEthereumSleepTime       = 4500     //ms
	DefaultTendermintStartHeight   = int64(0) // 0 will not change the db at start
	DefaultEthereumStartHeight     = int64(0) // 0 will not change the db at start
	DefaultPBridgeHome             = os.ExpandEnv("$HOME/.persistenceBridge")
	DefaultDenom                   = "stake"
	DefaultEthGasLimit             = uint64(500000)
	DefaultBroadcastMode           = flags.BroadcastAsync
	DefaultMinimumWrapAmount       = int64(5000000)
	DefaultAccountPrefix           = "cosmos"
	DefaultTendermintNode          = "http://127.0.0.1:26657"
	DefaultTendermintChainId       = "test"
	DefaultTendermintCoinType      = uint32(118)
	DefaultKafkaPorts              = "localhost:9092"
	DefaultCASPWaitTime            = 8 * time.Second
	DefaultCASPMaxAttempts         = 5
	DefaultRPCEndpoint             = "localhost:4040"
	DefaultTendermintMaxTxAttempts = 5
)
