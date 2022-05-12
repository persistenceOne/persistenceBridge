/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package casp

type UncompressedPublicKeysResponse struct {
	TotalItems   int64    `json:"totalItems"`
	Items        []string `json:"items"`
	Chains       []string `json:"chains"`
	AccountName  string   `json:"accountName"`
	AccountIndex uint32   `json:"accountIndex"`
}
