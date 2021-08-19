package db

import (
	"bytes"
	"encoding/binary"
	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/outgoingTx"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateStoreKey(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	pStakeConfig := configuration.InitConfig()
	_, err = toml.DecodeFile(filepath.Join(dirname, "/.persistenceBridge/config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}

	Txhash := common.BytesToHash([]byte("0x800b423ae1dfaf59de9e31fa45ebe0f57268949a8849fc2bd5f054b7c40eb18a"))
	Address := common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa"))
	amt := new(big.Int)
	amt.SetInt64(1000)

	wrapTokenMsg := outgoingTx.WrapTokenMsg{
		Address: Address,
		Amount:  amt,
	}
	tx := []outgoingTx.WrapTokenMsg{wrapTokenMsg}

	ethTransaction := EthereumBroadcastedWrapTokenTransaction{
		TxHash:   Txhash,
		Messages: tx,
	}
	err = SetEthereumTx(ethTransaction)
	if err != nil {
		t.Fatalf("Error %v", err.Error())
	}

	key := ethTransaction.TxHash.Bytes()
	storeKey := EthereumBroadcastedWrapTokenTransactionPrefix.GenerateStoreKey(key)
	Bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(Bytes, uint16(EthereumBroadcastedWrapTokenTransactionPrefix))

	//var newStoreKey []byte
	Key := append(Bytes, key...)
	res := bytes.Compare(storeKey, Key)
	if res != 0 {
		t.Fatal("Error")
	}
	db.Close()
}
