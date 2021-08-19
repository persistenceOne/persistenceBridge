package db

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func Test_GetUnboundEpochTimeandgetUnboundEpochTime(t *testing.T) {

	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	var epochTime int64 = 4772132
	err = SetUnboundEpochTime(epochTime)
	if err != nil {
		t.Fatalf("%v", err.Error())
	}

	var u UnboundEpochTime
	key := unboundEpochTimePrefix.GenerateStoreKey([]byte(unboundEpochTime))
	b, err := get(key)
	if err != nil {
		t.Fatalf("%v", err.Error())
	}
	err = json.Unmarshal(b, &u)
	if err != nil {
		t.Fatalf("%v", err.Error())
	}

	newUnboundEpochTime, err := GetUnboundEpochTime()
	if err != nil {
		t.Fatalf("%v", err.Error())
	}

	if newUnboundEpochTime.Epoch != epochTime {
		t.Fatalf("Could not get the correct Epoch Time, got %v expected %v", newUnboundEpochTime, epochTime)
	}
	db.Close()

}

func TestUnboundEpochTime_Key(t *testing.T) {

	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	if err != nil {
		t.Fatalf("%v", err.Error())
	}

	var epochTime int64 = 4772132
	unboundEpochTime := UnboundEpochTime{
		Epoch: epochTime,
	}

	Key := unboundEpochTime.Key()
	newKey := unboundEpochTime.prefix().GenerateStoreKey([]byte("UNBOUND_EPOCH_TIME"))
	if bytes.Compare(Key, newKey) != 0 {
		t.Fatalf("Wrong key returned,expected %v got %v", Key, newKey)
	}
	db.Close()

}

func TestUnboundEpochTime_Value(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	if err != nil {
		t.Fatalf("%v", err.Error())
	}

	var epochTime int64 = 4772132
	unboundEpochTime := UnboundEpochTime{
		Epoch: epochTime,
	}
	Value, err := json.Marshal(unboundEpochTime)
	if err != nil {
		t.Fatalf("%v", err.Error())
	}
	newValue, err := unboundEpochTime.Value()
	if err != nil {
		t.Fatalf("%v", err.Error())
	}
	if bytes.Compare(Value, newValue) != 0 {
		t.Fatalf("Wrong value returned,expected %v got %v ", Value, newValue)
	}
	db.Close()

}

func TestUnboundEpochTime_prefix(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	if err != nil {
		t.Fatalf("%v", err.Error())
	}

	var epochTime int64 = 4772132
	unboundEpochTime := UnboundEpochTime{
		Epoch: epochTime,
	}

	Prefix := unboundEpochTime.prefix()

	if Prefix != unboundEpochTimePrefix {
		t.Fatalf("Wrong prefix returned,expected %v got %v ", unboundEpochTimePrefix, Prefix)
	}
	db.Close()
}
