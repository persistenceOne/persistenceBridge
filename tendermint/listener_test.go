package tendermint

import (
	"github.com/BurntSushi/toml"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestStartListening(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	fileName := strings.Join([]string{dirname,"/.persistenceBridge/chain.json"},"")
	initAndStartChain, errInit := InitializeAndStartChain(fileName, "336h",dirname)
	pStakeConfig := configuration.InitConfig()
	_, err := toml.DecodeFile(filepath.Join(dirname, "/.persistenceBridge/config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}
	if errInit != nil {
		t.Errorf("Error Initializing Chain: %v",errInit)
	}
	initClientCtx := client.Context{}.
		WithHomeDir(constants.DefaultPBridgeHome)
	protoCodec := codec.NewProtoCodec(initClientCtx.InterfaceRegistry)
	StartListening(initClientCtx.WithHomeDir(dirname), initAndStartChain, pStakeConfig.Kafka.Brokers, protoCodec, time.Duration(200)*time.Millisecond)
}

func Test_getAllTxResults(t *testing.T) {

}
