package outgoingTx

import (
	"github.com/BurntSushi/toml"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/relayer/relayer"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/stretchr/testify/require"
	"log"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLogMessagesAndBroadcast(t *testing.T) {
	type args struct {
		chain         *relayer.Chain
		msgs          []sdk.Msg
		timeoutHeight uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *sdk.TxResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LogMessagesAndBroadcast(tt.args.chain, tt.args.msgs, tt.args.timeoutHeight)
			if (err != nil) != tt.wantErr {
				t.Errorf("LogMessagesAndBroadcast() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LogMessagesAndBroadcast() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_broadcastTMTx(t *testing.T) {
	type args struct {
		chain         *relayer.Chain
		fromPublicKey types.PubKey
		sigBytes      []byte
		txBuilder     client.TxBuilder
		txFactory     tx.Factory
	}
	tests := []struct {
		name    string
		args    args
		want    *sdk.TxResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := broadcastTMTx(tt.args.chain, tt.args.fromPublicKey, tt.args.sigBytes, tt.args.txBuilder, tt.args.txFactory)
			if (err != nil) != tt.wantErr {
				t.Errorf("broadcastTMTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("broadcastTMTx() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTMBytesToSign(t *testing.T) {
	type args struct {
		chain         *relayer.Chain
		fromPublicKey types.PubKey
		msgs          []sdk.Msg
		memo          string
		timeoutHeight uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		want1   client.TxBuilder
		want2   tx.Factory
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := getTMBytesToSign(tt.args.chain, tt.args.fromPublicKey, tt.args.msgs, tt.args.memo, tt.args.timeoutHeight)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTMBytesToSign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTMBytesToSign() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getTMBytesToSign() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("getTMBytesToSign() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_getTMSignature(t *testing.T) {
	pStakeConfig := configuration.InitConfig()
	_, err := toml.DecodeFile(filepath.Join("/Users/ankitkumar/.persistenceBridge/", "config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}
	dataToSign := []string{"55C53F5D490297900CEFA825D0C8E8E9532EE8A118ABE7D8570762CD38BE9818"}
	bytesToSign :=[]byte(strings.Join(dataToSign,""))
	tmSignature, err := getTMSignature(bytesToSign)
	if err != nil {
		t.Errorf("Error getting TM signature: \n %v",err)
	}
	require.NotNil(t, tmSignature)
	require.Equal(t, 64, len(tmSignature))
	//fmt.Println(tmSignature)

}

func Test_setTMPublicKey(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setTMPublicKey(); (err != nil) != tt.wantErr {
				t.Errorf("setTMPublicKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_tendermintSignAndBroadcastMsgs(t *testing.T) {
	type args struct {
		chain         *relayer.Chain
		msgs          []sdk.Msg
		memo          string
		timeoutHeight uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *sdk.TxResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tendermintSignAndBroadcastMsgs(tt.args.chain, tt.args.msgs, tt.args.memo, tt.args.timeoutHeight)
			if (err != nil) != tt.wantErr {
				t.Errorf("tendermintSignAndBroadcastMsgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tendermintSignAndBroadcastMsgs() got = %v, want %v", got, tt.want)
			}
		})
	}
}
