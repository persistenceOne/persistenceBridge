package db

import (
	"encoding/json"
)

const unboundEpochTime = "UNBOUND_EPOCH_TIME"

type UnboundEpochTime struct {
	Epoch int64
}

var _ DBI = &UnboundEpochTime{}

func (u *UnboundEpochTime) prefix() storeKeyPrefix {
	return unboundEpochTimePrefix
}

func (u *UnboundEpochTime) Key() []byte {
	return u.prefix().GenerateStoreKey([]byte(unboundEpochTime))
}

func (u *UnboundEpochTime) Value() ([]byte, error) {
	return json.Marshal(*u)
}

func (u *UnboundEpochTime) Validate() error {
	return nil
}

func GetUnboundEpochTime() (UnboundEpochTime, error) {
	var u UnboundEpochTime
	key := unboundEpochTimePrefix.GenerateStoreKey([]byte(unboundEpochTime))
	b, err := get(key)
	if err != nil {
		return u, err
	}
	err = json.Unmarshal(b, &u)
	return u, err
}

func SetUnboundEpochTime(epochTime int64) error {
	u := UnboundEpochTime{
		Epoch: epochTime,
	}
	return set(&u)
}
