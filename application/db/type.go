/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

type DBI interface {
	prefix() storeKeyPrefix
	Key() []byte
	Value() ([]byte, error)
	Validate() error
}
