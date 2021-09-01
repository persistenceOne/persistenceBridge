package tendermint

import (
	"github.com/BurntSushi/toml"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func Test_beta(t *testing.T) {
	configuration.InitConfig()
	appConfig := test.GetCmdWithConfig()
	configuration.SetConfig(&appConfig)
	
	limiter := db.AccountLimiter{AccountAddress: sdk.AccAddress("cosmos1t48p2wwqafhsgmf0uf7wcmk3zkq9f5d7lzl74n"), Amount: sdk.NewInt(int64(1))}
	sendAmount, refundAmt := beta(limiter,sdk.NewInt(int64(200)))
	require.Equal(t,reflect.TypeOf(sdk.Int{}) ,reflect.TypeOf(sendAmount))
	require.Equal(t,reflect.TypeOf(sdk.Int{}) ,reflect.TypeOf(refundAmt))
}

func Test_collectAllWrapAndRevertTxs(t *testing.T) {

	//tmwrap := collectAllWrapAndRevertTxs(client.Context{},&tmCoreTypes.ResultTx{}),
}

func Test_getMaxLimit(t *testing.T) {
	response := getMaxLimit()
	require.LessOrEqual(t, 50000,response)
	require.GreaterOrEqual(t,2147483646, response )
	require.Equal(t, reflect.TypeOf(1),reflect.TypeOf(response))
}

func Test_handleTxSearchResult(t *testing.T) {
	//kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec
}

func Test_revertCoins(t *testing.T) {
	//toAddress string, coins sdk.Coins, kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec
}

func Test_wrapOrRevert(t *testing.T) {
	//tmWrapOrReverts []tmWrapOrRevert, kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec

}
