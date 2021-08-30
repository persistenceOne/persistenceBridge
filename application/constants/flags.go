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
	FlagTendermintSleepTime   = "tmSleepTime"
	FlagEthereumSleepTime     = "ethSleepTime"
	FlagTendermintStartHeight = "tmStart"
	FlagEthereumStartHeight   = "ethStart"
	FlagDenom                 = "denom"
	FlagEthGasLimit           = "ethGasLimit"
	FlagKafkaPorts            = "kafkaPorts"
	FlagBroadcastMode         = "tmBroadcastMode"
	FlagCASPURL               = "caspURL"
	FlagCASPVaultID           = "caspVaultID"
	FlagCASPTMPublicKey       = "caspTMPublicKey"
	FlagCASPEthPublicKey      = "caspEthPublicKey"
	FlagCASPSignatureWaitTime = "caspSignatureWaitTime"
	FlagCASPApiToken          = "caspApiToken"
	FlagCASPConcurrentKey     = "caspConcurrentKeyUsage"
	FlagRPCEndpoint           = "rpc-endpoint"
	FlagMinimumWrapAmount     = "minWrapAmt"
	FlagTelegramBotToken      = "botToken"
	FlagTelegramChatID        = "botChatID"
	FlagShowDebugLog          = "debug"
)

var (
	DefaultTimeout               = "10s"
	DefaultEthereumEndPoint      = "wss://ropsten.infura.io/ws/v3/b21966541db246d398fb31402eec2c14"
	DefaultTendermintSleepTime   = 3000     //ms
	DefaultEthereumSleepTime     = 4500     //ms
	DefaultTendermintStartHeight = int64(0) // 0 will not change the db at start
	DefaultEthereumStartHeight   = int64(0) // 0 will not change the db at start
	DefaultPBridgeHome           = os.ExpandEnv("$HOME/.persistenceBridge")
	DefaultDenom                 = "uatom"
	DefaultEthGasLimit           = uint64(500000)
	DefaultBroadcastMode         = flags.BroadcastAsync
	DefaultMinimumWrapAmount     = int64(5000000)
	DefaultKafkaPorts            = "localhost:9092"
	DefaultCASPSignatureWaitTime = 8 * time.Second
	DefaultRPCEndpoint           = "localhost:4040"
)
