package constants

import (
	"os"
	"time"
)

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
	FlagKafkaPorts            = "kafkaPorts"
	FlagBroadcastMode         = "tmBroadcastMode"
	FlagCASPURL               = "caspURL"
	FlagCASPVaultID           = "caspVaultID"
	FlagCASPPublicKey         = "caspPublicKey"
	FlagCASPSignatureWaitTime = "caspSignatureWaitTime"
	FlagCASPApiToken          = "caspApiToken"
	FlagCASPCoin              = "caspCoin"
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
	DefaultBroadcastMode         = "sync"
	DefaultCASPUrl               = "https://65.2.149.241:443"
	DefaultCASPVaultID           = "509fd89a-762a-40ec-bd4b-0745b06e2d3d"
	DefaultCASPAPI               = "Bearer cHVuZWV0TmV3QXBpa2V5MTI6OWM1NDBhMzAtNTQ5NC00ZDdhLTljODktODA3MDZiNWNhYzQ1"
	DefaultCASPPublicKey         = "3056301006072A8648CE3D020106052B8104000A034200044F717AE01D84C0827054A4505D779632072F923C811B8A2A2D12B4A55A1B59A4DB2F5FEF4B52B7D4DD08B8047B4ACD565488EAA88CDC2A99EE1E796AD7D1BDDA"
	DefaultCASPCoin              = DefaultCoinType
	DefaultKafkaPorts            = "localhost:9092"
	DefaultCASPSignatureWaitTime = 8 * time.Second
)
