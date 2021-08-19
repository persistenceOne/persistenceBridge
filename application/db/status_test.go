package db

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func Test_getStatusandSetStatus(t *testing.T) {

	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")

	Name := "tx1"
	var LastCheckHeight int64 = 4772132
	err = setStatus(Name, LastCheckHeight)
	if err != nil {
		t.Fatalf("%v", err.Error())
	}

	var status Status
	status.Name = Name
	b, err := get(status.Key())
	if err != nil {
		t.Fatalf("%v", err.Error())
	}
	err = json.Unmarshal(b, &status)

	newStatus, err := getStatus(Name)
	if err != nil {
		t.Fatalf("%v", err.Error())
	}

	if reflect.DeepEqual(newStatus, status) != true {
		t.Fatalf("Could not get the correct status, expected %v got %v", status, newStatus)
	}
	db.Close()

}

func TestStatus_Key(t *testing.T) {

	Name := "tx1"
	var LastCheckHeight int64 = 4772132

	status := Status{
		Name,
		LastCheckHeight,
	}
	Key := status.prefix().GenerateStoreKey([]byte(status.Name))
	newKey := status.Key()
	if bytes.Compare(Key, newKey) != 0 {
		t.Fatalf("Wrong key returned,expected %v got %v", Key, newKey)
	}
}
func TestStatus_Value(t *testing.T) {
	Name := "tx1"
	var LastCheckHeight int64 = 4772132

	status := Status{
		Name,
		LastCheckHeight,
	}
	Value, err := json.Marshal(status)
	if err != nil {
		t.Fatalf("%v", err.Error())
	}
	newValue, err := status.Value()
	if err != nil {
		t.Fatalf("%v", err.Error())
	}
	if bytes.Compare(Value, newValue) != 0 {
		t.Fatalf("Wrong value returned,expected %v got %v ", Value, newValue)
	}
}

func TestStatus_prefix(t *testing.T) {
	Name := "tx1"
	var LastCheckHeight int64 = 4772132

	status := Status{
		Name,
		LastCheckHeight,
	}
	Prefix := status.prefix()
	if Prefix != statusPrefix {
		t.Fatalf("Wrong prefix returned,expected %v got %v ", statusPrefix, Prefix)
	}

}
