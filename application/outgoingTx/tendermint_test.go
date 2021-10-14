package outgoingTx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/relayer/helpers"
	"github.com/cosmos/relayer/relayer"
	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	caspQueries "github.com/persistenceOne/persistenceBridge/application/rest/casp"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestLogMessagesAndBroadcast(t *testing.T) {
	test.SetTestConfig()
	tenderMintAddress, err := casp.GetTendermintAddress()
	require.Equal(t, nil, err)
	ethAddress, err := casp.GetEthAddress()
	require.Equal(t, nil, err)
	configuration.SetCASPAddresses(tenderMintAddress, ethAddress)
	chain := setUpChain(t)

	msg := &bankTypes.MsgSend{
		FromAddress: configuration.GetAppConfig().Tendermint.GetWrapAddress(),
		ToAddress:   "cosmos19u3y3gx35509fwxj5s0fzsz85qs452d8t4da06",
		Amount:      sdk.NewCoins(sdk.NewInt64Coin(configuration.GetAppConfig().Tendermint.Denom, 1)),
	}
	txResponse, err := LogMessagesAndBroadcast(chain, []sdk.Msg{msg}, 200)
	require.Equal(t, nil, err)

	re := regexp.MustCompile(`^[0-9a-fA-f]{64}`)
	require.NotNil(t, txResponse)
	require.Equal(t, true, re.MatchString(txResponse.TxHash))
	require.Equal(t, reflect.TypeOf(&sdk.TxResponse{}), reflect.TypeOf(txResponse))
	require.Equal(t, reflect.TypeOf(""), reflect.TypeOf(txResponse.String()))
}

func TestBroadcastTMTx(t *testing.T) {
	test.SetTestConfig()
	uncompressedPublicKeys, err := caspQueries.GetUncompressedTMPublicKeys()
	require.Equal(t, nil, err)
	tmpPubKey := casp.GetTMPubKey(uncompressedPublicKeys.Items[0])
	tmAddress, err := casp.GetTendermintAddress()
	require.Equal(t, nil, err)
	ethAddress, err := casp.GetEthAddress()
	require.Equal(t, nil, err)
	configuration.SetCASPAddresses(tmAddress, ethAddress)
	chain := setUpChain(t)
	msg := &bankTypes.MsgSend{
		FromAddress: configuration.GetAppConfig().Tendermint.GetWrapAddress(),
		ToAddress:   "cosmos19u3y3gx35509fwxj5s0fzsz85qs452d8t4da06",
		Amount:      sdk.NewCoins(sdk.NewInt64Coin(configuration.GetAppConfig().Tendermint.Denom, 1)),
	}

	bytesToSign, txB, txF, err := getTMBytesToSign(chain, tmpPubKey, []sdk.Msg{msg}, "pStake@PersistenceOne", 200)
	require.Equal(t, nil, err)
	signature, err := getTMSignature(bytesToSign)
	require.Equal(t, nil, err)
	broadcastTMmsg, err := broadcastTMTx(chain, tmpPubKey, signature, txB, txF)
	require.Equal(t, nil, err)
	re := regexp.MustCompile(`^[0-9a-fA-F]`)
	require.Equal(t, true, re.MatchString(broadcastTMmsg.TxHash))
	require.NotNil(t, broadcastTMmsg)
	require.Equal(t, reflect.TypeOf(&sdk.TxResponse{}), reflect.TypeOf(broadcastTMmsg))
}

func TestGetTMBytesToSign(t *testing.T) {
	test.SetTestConfig()
	uncompressedPublicKeys, err := caspQueries.GetUncompressedTMPublicKeys()
	require.Equal(t, nil, err)
	tmpPubKey := casp.GetTMPubKey(uncompressedPublicKeys.Items[0])
	tmAddress, err := casp.GetTendermintAddress()
	require.Equal(t, nil, err)
	ethAddress, err := casp.GetEthAddress()
	require.Equal(t, nil, err)
	configuration.SetCASPAddresses(tmAddress, ethAddress)
	chain := setUpChain(t)

	msg := &bankTypes.MsgSend{
		FromAddress: configuration.GetAppConfig().Tendermint.GetWrapAddress(),
		ToAddress:   "cosmos19u3y3gx35509fwxj5s0fzsz85qs452d8t4da06",
		Amount:      sdk.NewCoins(sdk.NewInt64Coin(configuration.GetAppConfig().Tendermint.Denom, 1)),
	}

	tmBytesSignBytes, txBuilder, txFactory, errorGettingTMBytes := getTMBytesToSign(chain, tmpPubKey, []sdk.Msg{msg}, "pStake@PersistenceOne", 200)
	if errorGettingTMBytes != nil {
		t.Errorf("Error Getting TM Bytes to Sign: %v", errorGettingTMBytes)
	}
	require.Equal(t, "pStake@PersistenceOne", txFactory.Memo())
	require.NotNil(t, tmBytesSignBytes)
	require.NotNil(t, txBuilder)
	require.NotNil(t, txFactory)
}

func TestGetTMSignature(t *testing.T) {
	test.SetTestConfig()

	dataToSign := []string{"55C53F5D490297900CEFA825D0C8E8E9532EE8A118ABE7D8570762CD38BE9818"}
	bytesToSign := []byte(strings.Join(dataToSign, ""))
	tmSignature, err := getTMSignature(bytesToSign)
	if err != nil {
		t.Errorf("Error getting TM signature: \n %v", err)
	}
	require.Equal(t, 64, len(tmSignature))
	require.Equal(t, reflect.TypeOf([]byte{}), reflect.TypeOf(tmSignature))
	require.NotNil(t, tmSignature)
	require.Equal(t, 64, len(tmSignature))
}

func TestSetTMPublicKey(t *testing.T) {
	test.SetTestConfig()
	err := setTMPublicKey()
	require.Equal(t, nil, err)
	require.NotNil(t, tmPublicKey)
	require.Equal(t, 20, len(tmPublicKey.Address()))
}

func TestTendermintSignAndBroadcastMsgs(t *testing.T) {
	test.SetTestConfig()
	tmAddress, err := casp.GetTendermintAddress()
	require.Equal(t, nil, err)
	ethAddress, err := casp.GetEthAddress()
	require.Equal(t, nil, err)
	configuration.SetCASPAddresses(tmAddress, ethAddress)

	msg := &bankTypes.MsgSend{
		FromAddress: configuration.GetAppConfig().Tendermint.GetWrapAddress(),
		ToAddress:   "cosmos19u3y3gx35509fwxj5s0fzsz85qs452d8t4da06",
		Amount:      sdk.NewCoins(sdk.NewInt64Coin(configuration.GetAppConfig().Tendermint.Denom, 1)),
	}

	chain := setUpChain(t)
	tmSignAndBroadcastMsg, err := tendermintSignAndBroadcastMsgs(chain, []sdk.Msg{msg}, "", 0)
	require.Equal(t, nil, err)
	re := regexp.MustCompile(`^[0-9a-fA-F]{64}`)
	require.Equal(t, true, re.MatchString(tmSignAndBroadcastMsg.TxHash))
	require.Equal(t, reflect.TypeOf(&sdk.TxResponse{}), reflect.TypeOf(tmSignAndBroadcastMsg))
	require.NotNil(t, tmSignAndBroadcastMsg)
}

func setUpChain(t *testing.T) *relayer.Chain {
	dirname, _ := os.UserHomeDir()
	homePath := strings.Join([]string{dirname, "/.persistenceBridge"}, "/")

	chain := &relayer.Chain{}
	chain.Key = "unusedKey"
	chain.ChainID = configuration.GetAppConfig().Tendermint.ChainID
	chain.RPCAddr = configuration.GetAppConfig().Tendermint.Node
	chain.AccountPrefix = configuration.GetAppConfig().Tendermint.AccountPrefix
	chain.GasAdjustment = configuration.GetAppConfig().Tendermint.GasAdjustment
	chain.GasPrices = configuration.GetAppConfig().Tendermint.GasPrice + configuration.GetAppConfig().Tendermint.Denom
	chain.TrustingPeriod = "21h"

	to, err := time.ParseDuration("10s")
	require.Equal(t, nil, err)

	err = chain.Init(homePath, to, nil, true)
	require.Equal(t, nil, err)

	if chain.KeyExists(chain.Key) {
		err = chain.Keybase.Delete(chain.Key)
		require.Equal(t, nil, err)
	}

	_, err = helpers.KeyAddOrRestore(chain, chain.Key, uint32(118))
	require.Equal(t, nil, err)

	err = chain.Start()
	require.Equal(t, nil, err)

	return chain
}
