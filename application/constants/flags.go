package constants

import (
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/ethereum/go-ethereum/common"
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
	FlagRPCEndpoint           = "rpc"
)

var (
	DefaultTimeout                 = "10s"
	DefaultEthereumEndPoint        = "wss://ropsten.infura.io/ws/v3/b21966541db246d398fb31402eec2c14"
	DefaultTendermintSleepTime     = 3000     //ms
	DefaultEthereumSleepTime       = 4500     //ms
	DefaultTendermintStartHeight   = int64(0) // 0 will not change the db at start
	DefaultEthereumStartHeight     = int64(0) // 0 will not change the db at start
	DefaultPBridgeHome             = os.ExpandEnv("$HOME/.persistenceBridge")
	DefaultDenom                   = "uatom"
	DefaultPStakeAddress           = "cosmos15vs9hfghf3xpsqshw98gq6mtt55wmhlgxf83pd"                   //TODO should be derived not given
	DefaultBridgeAdmin             = common.HexToAddress("0xfCd7b44E0F250928aEC442ebc5E7bc0e4B38a8D5") //TODO should be derived not given
	DefaultEthGasLimit             = uint64(500000)
	DefaultBroadcastMode           = flags.BroadcastAsync
	DefaultCASPUrl                 = "https://65.2.149.241:443"
	DefaultCASPVaultID             = "509fd89a-762a-40ec-bd4b-0745b06e2d3d"
	DefaultCASPAPI                 = "Bearer cHVuZWV0TmV3QXBpa2V5MTI6OWM1NDBhMzAtNTQ5NC00ZDdhLTljODktODA3MDZiNWNhYzQ1"
	DefaultCASPTendermintPublicKey = "3056301006072A8648CE3D020106052B8104000A034200044F717AE01D84C0827054A4505D779632072F923C811B8A2A2D12B4A55A1B59A4DB2F5FEF4B52B7D4DD08B8047B4ACD565488EAA88CDC2A99EE1E796AD7D1BDDA"
	DefaultCASPEthereumPublicKey   = "3056301006072A8648CE3D020106052B8104000A03420004B40777F842A9F8BB7ECB94785926D725EB1F96611DC2B2C424EBC8BD1A9B7651302DC7A55301D560D599B3F72D630353325FAED84514C4ECD58330B965A1946A"
	DefaultKafkaPorts              = "localhost:9092"
	DefaultCASPSignatureWaitTime   = 8 * time.Second
	DefaultRPCEndpoint             = "localhost:4040"
)
