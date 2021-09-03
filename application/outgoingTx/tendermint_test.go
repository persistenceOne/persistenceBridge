package outgoingTx

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/relayer/helpers"
	"github.com/cosmos/relayer/relayer"
	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	caspQueries "github.com/persistenceOne/persistenceBridge/application/rest/casp"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestLogMessagesAndBroadcast(t *testing.T) {
	configuration.InitConfig()
	appConfig := test.GetCmdWithConfig()
	configuration.SetConfig(&appConfig)
	dirname, _ := os.UserHomeDir()
	tenderMintAddress, errorTm := casp.GetTendermintAddress()
	require.Nil(t, errorTm,"Error getting Tendermint address")
	configuration.SetPStakeAddress(tenderMintAddress)
	chain := &relayer.Chain{}
	byte,err := ioutil.ReadFile(strings.Join([]string{dirname,"/.persistenceBridge/chain.json"},""))
	require.Nil(t, err,"No config files found")
	json.Unmarshal(byte, chain)
	to, err := time.ParseDuration("200")
	err = chain.Init(dirname, to, nil, true)
	if err != nil {
		return
	}
	if chain.KeyExists(chain.Key) {
		logging.Info("deleting old key", chain.Key)
		err = chain.Keybase.Delete(chain.Key)
		if err != nil {
			return
		}
	}
	ko, err := helpers.KeyAddOrRestore(chain, chain.Key, uint32(118))
	if err != nil {
		return
	}

	logging.Info("Relayer Chain Keys added [NOT TO BE USED]:", ko.Address)
	msg := &bankTypes.MsgSend{
		FromAddress: configuration.GetAppConfig().Tendermint.GetPStakeAddress(),
		ToAddress:   "cosmos19u3y3gx35509fwxj5s0fzsz85qs452d8t4da06",
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("validatortoken",1)),
	}
	msgs := []sdk.Msg{msg}
	loggedmessage, errr := LogMessagesAndBroadcast(chain ,msgs,200)
	if errr != nil {
		t.Errorf("Error logging messaged" +
			": %v",errr)
	}
	re := regexp.MustCompile(`^[0-9a-fA-f]{64}`)
	require.NotNil(t, loggedmessage)
	require.Equal(t, true,re.MatchString(loggedmessage.TxHash))
	require.Equal(t, reflect.TypeOf(&sdk.TxResponse{}),reflect.TypeOf(loggedmessage))
	require.Equal(t, reflect.TypeOf(""),reflect.TypeOf(loggedmessage.String()))
}

func TestBroadcastTMTx(t *testing.T) {
	configuration.InitConfig()
	appConfig := test.GetCmdWithConfig()
	configuration.SetConfig(&appConfig)
	dirname, _ := os.UserHomeDir()
	uncompressedPublicKeys, err := caspQueries.GetUncompressedTMPublicKeys()
	tmpPubKey := casp.GetTMPubKey(uncompressedPublicKeys.PublicKeys[0])
	tmAddress, err := casp.GetTendermintAddress()
	configuration.SetPStakeAddress(tmAddress)
	chain := &relayer.Chain{}
	byte,err := ioutil.ReadFile(strings.Join([]string{dirname,"/.persistenceBridge/chain.json"},""))
	require.Nil(t, err,"No config files found")
	json.Unmarshal(byte, chain)
	to, err := time.ParseDuration("200")
	err = chain.Init(dirname, to, nil, true)
	if err != nil {
		return
	}
	if chain.KeyExists(chain.Key) {
		logging.Info("deleting old key", chain.Key)
		err = chain.Keybase.Delete(chain.Key)
		if err != nil {
			return
		}
	}
	_, erroKey := helpers.KeyAddOrRestore(chain, chain.Key, uint32(118))
	require.Nil(t, erroKey,"Key error!")
	msg := &bankTypes.MsgSend{
		FromAddress: configuration.GetAppConfig().Tendermint.GetPStakeAddress(),
		ToAddress:   "cosmos19u3y3gx35509fwxj5s0fzsz85qs452d8t4da06",
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("validatortoken",1)),
	}
	msgs := []sdk.Msg{msg}

	tendermintPublicKeyError := setTMPublicKey()
	if tendermintPublicKeyError != nil {
		t.Errorf("Error setting tenderMintpublic Key: %v",tendermintPublicKeyError)
	}
	bytesToSign, txB, txF, errTmBytesToSign := getTMBytesToSign(chain,tmpPubKey,msgs,"pStake@PersistenceOne",200)
	if errTmBytesToSign != nil {
		t.Errorf("Error Signing TM bytes: %v", errTmBytesToSign)
	}
	signature, errTmSign := getTMSignature(bytesToSign)
	if errTmSign != nil {
		t.Errorf("Error getting TM sign: %v",errTmSign)
	}
	broadcastTMmsg, errBroadcastTMTx := broadcastTMTx(chain,tmpPubKey,signature,txB,txF)
	if errBroadcastTMTx != nil {
		t.Errorf("Error Broadcasting TM sign: %v",errBroadcastTMTx)
	}
	re := regexp.MustCompile(`^[0-9a-fA-F]`)
	require.Equal(t, true,re.MatchString(broadcastTMmsg.TxHash))
	require.Equal(t, 66, broadcastTMmsg.Size())
	require.NotNil(t, broadcastTMmsg)
	require.Equal(t, reflect.TypeOf(&sdk.TxResponse{}),reflect.TypeOf(broadcastTMmsg))
}

func TestGetTMBytesToSign(t *testing.T) {
	configuration.InitConfig()
	appConfig := test.GetCmdWithConfig()
	configuration.SetConfig(&appConfig)
	dirname, _ := os.UserHomeDir()
	uncompressedPublicKeys, err := caspQueries.GetUncompressedTMPublicKeys()
	tmpPubKey := casp.GetTMPubKey(uncompressedPublicKeys.PublicKeys[0])
	tmAddress, err := casp.GetTendermintAddress()
	configuration.SetPStakeAddress(tmAddress)
	chain := &relayer.Chain{}
	byte,err := ioutil.ReadFile(strings.Join([]string{dirname,"/.persistenceBridge/chain.json"},""))
	require.Nil(t, err,"No config files found")
	json.Unmarshal(byte, chain)
	to, err := time.ParseDuration("200")
	err = chain.Init(dirname, to, nil, true)
	if err != nil {
		return
	}
	if chain.KeyExists(chain.Key) {
		logging.Info("deleting old key", chain.Key)
		err = chain.Keybase.Delete(chain.Key)
		if err != nil {
			return
		}
	}
	_, erroKey := helpers.KeyAddOrRestore(chain, chain.Key, uint32(118))
	require.Nil(t, erroKey,"Key error!")

	msg := &bankTypes.MsgSend{
		FromAddress: configuration.GetAppConfig().Tendermint.GetPStakeAddress(),
		ToAddress:   "cosmos19u3y3gx35509fwxj5s0fzsz85qs452d8t4da06",
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("validatortoken",1)),
	}
	msgs := []sdk.Msg{msg}

	tendermintPublicKeyError := setTMPublicKey()
	if tendermintPublicKeyError != nil {
		t.Errorf("Error setting tenderMintpublic Key: %v",tendermintPublicKeyError)
	}
	tmBytesSignBytes, txBuilder, txFactory, errorGettingTMBytes := getTMBytesToSign(chain,tmpPubKey,msgs,"pStake@PersistenceOne",200)
	if errorGettingTMBytes != nil {
		t.Errorf("Error Getting TM Bytes to Sign: %v",errorGettingTMBytes)
	}
	require.Equal(t, reflect.TypeOf(byte),reflect.TypeOf(tmBytesSignBytes))
	require.Equal(t, 290, len(tmBytesSignBytes))
	require.Equal(t, "pStake@PersistenceOne", txFactory.Memo())
	require.NotNil(t, tmBytesSignBytes)
	require.NotNil(t, txBuilder)
	require.NotNil(t, txFactory)
}

func TestGetTMSignature(t *testing.T) {
	configuration.InitConfig()
	appConfig := test.GetCmdWithConfig()
	configuration.SetConfig(&appConfig)
	
	dataToSign := []string{"55C53F5D490297900CEFA825D0C8E8E9532EE8A118ABE7D8570762CD38BE9818"}
	bytesToSign :=[]byte(strings.Join(dataToSign,""))
	tmSignature, err := getTMSignature(bytesToSign)
	if err != nil {
		t.Errorf("Error getting TM signature: \n %v",err)
	}
	require.Equal(t, 64,len(tmSignature))
	require.Equal(t, reflect.TypeOf([]byte{}),reflect.TypeOf(tmSignature))
	require.NotNil(t, tmSignature)
	require.Equal(t, 64, len(tmSignature))
}

func TestSetTMPublicKey(t *testing.T) {
	configuration.InitConfig()
	appConfig := test.GetCmdWithConfig()
	configuration.SetConfig(&appConfig)
	
	tendermintPublicKey := setTMPublicKey()
	if tendermintPublicKey != nil {
		t.Errorf("Error setting Tendermint publickey: %v",tendermintPublicKey)
	}
	require.NotNil(t, tmPublicKey)
	require.Equal(t, 20,len(tmPublicKey.Address()))
}

func TestTendermintSignAndBroadcastMsgs(t *testing.T) {
	configuration.InitConfig()
	appConfig := test.GetCmdWithConfig()
	configuration.SetConfig(&appConfig)
	dirname, _ := os.UserHomeDir()
	tmAddress, err := casp.GetTendermintAddress()
	configuration.SetPStakeAddress(tmAddress)
	chain := &relayer.Chain{}
	//change file Name Type
	byte,err := ioutil.ReadFile(strings.Join([]string{dirname,"/.persistenceBridge/chain.json"},""))
	require.Nil(t, err,"No config files found")
	json.Unmarshal(byte, chain)
	to, err := time.ParseDuration("200")
	err = chain.Init(dirname, to, nil, true)
	if err != nil {
		return
	}
	if chain.KeyExists(chain.Key) {
		logging.Info("deleting old key", chain.Key)
		err = chain.Keybase.Delete(chain.Key)
		if err != nil {
			return
		}
	}
	_, erroKey := helpers.KeyAddOrRestore(chain, chain.Key, uint32(118))
	require.Nil(t, erroKey,"Key error!")

	msg := &bankTypes.MsgSend{
		FromAddress: configuration.GetAppConfig().Tendermint.GetPStakeAddress(),
		ToAddress:   "cosmos19u3y3gx35509fwxj5s0fzsz85qs452d8t4da06",
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("validatortoken",1)),
	}
	msgs := []sdk.Msg{msg}

	tendermintPublicKeyError := setTMPublicKey()
	if tendermintPublicKeyError != nil {
		t.Errorf("Error setting tenderMintpublic Key: %v",tendermintPublicKeyError)
	}
	tmSignAndBroadcastMsg, errSingAndBroadcast := tendermintSignAndBroadcastMsgs(chain,msgs,"",200)
	if errSingAndBroadcast != nil {
		t.Errorf("Error signing and Broadcasting msgs: %v",errSingAndBroadcast)
	}
	re := regexp.MustCompile(`^[0-9a-fA-F]{64}`)
	require.Equal(t, true,re.MatchString(tmSignAndBroadcastMsg.TxHash))
	require.Equal(t, reflect.TypeOf(&sdk.TxResponse{}),reflect.TypeOf(tmSignAndBroadcastMsg))
	require.NotNil(t, tmSignAndBroadcastMsg)
}
