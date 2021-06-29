package constants

import "os"

const (
	FlagTimeOut               = "timeout"
	FlagCoinType              = "coinType"
	FlagPBridgeHome           = "pBridgeHome"
	FlagEthereumEndPoint      = "ethEndPoint"
	FlagTendermintSleepTime   = "tmSleepTime"
	FlagEthereumSleepTime     = "ethSleepTime"
	FlagTendermintStartHeight = "tmStart"
	FlagEthereumStartHeight   = "ethStart"
	FlagDenom                 = "denom"
	FlagEthPrivateKey         = "ethPrivateKey"
	FlagEthGasLimit           = "ethGasLimit"
)

var (
	DefaultTimeout               = "10s"
	DefaultCoinType              = uint32(118)
	DefaultEthereumEndPoint      = "wss://ropsten.infura.io/ws/v3/b21966541db246d398fb31402eec2c14"
	DefaultTendermintSleepTime   = 3000     //ms
	DefaultEthereumSleepTime     = 4500     //ms
	DefaultTendermintStartHeight = int64(0) // 0 will not change the db at start
	DefaultEthereumStartHeight   = int64(0) // 0 will not change the db at start
	DefaultPBridgeHome           = os.ExpandEnv("$HOME/.persistenceBridge")
	DefaultDenom                 = "uatom"
	DefaultEthGasLimit           = uint64(3000000)
)
