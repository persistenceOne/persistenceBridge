package db

type DBI interface {
	prefix() storeKeyPrefix
	Key() []byte
	Value() ([]byte, error)
	Validate() error
}
