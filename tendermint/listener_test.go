package tendermint

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestStartListening(t *testing.T) {
//	homedir, _ := os.UserHomeDir()
//	fileName := strings.Join([]string{homedir,"/.persistenceBridge/chain.json"},"")
//	chain, _ := InitializeAndStartChain(fileName, "336h", homedir)
//	pStakeConfig := configuration.InitConfig()
//	_, err := toml.DecodeFile(filepath.Join(homedir, "/.persistenceBridge/config.toml"), &pStakeConfig)
//	if err != nil {
//		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
//	}
//	//ctx := client.Context{}
//	//EncodingConfig := app.MakeEncodingConfig()
//	//EncodingConfig := .MsgHandler{}
//	//protoCodec *codec.ProtoCodec
//	//StartListening(ctx , chain, pStakeConfig.Kafka.Brokers, EncodingConfig, time.Duration(200)*time.Millisecond)
}

func Test_getAllTxResults(t *testing.T) {
	homedir, _ := os.UserHomeDir()
	fileName := strings.Join([]string{homedir,"/.persistenceBridge/chain.json"},"")
	chain, _ := InitializeAndStartChain(fileName, "336h", homedir)
	ctx := context.Background()
	result, err := getAllTxResults(ctx, chain,0 )
	if err != nil {
		t.Errorf("Error getting all Tx Results: %v",err)
	}
	fmt.Println(result)
	require.Equal(t, reflect.TypeOf([]*coretypes.ResultTx{}) , reflect.TypeOf(result))
}
