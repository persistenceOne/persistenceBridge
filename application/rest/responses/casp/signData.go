/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package casp

import "fmt"

type PostSignDataResponse struct {
	OperationID string `json:"operationID"`
}

type ErrorResponse struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Details string `json:"details,omitempty"`
	Status  int64  `json:"status"`
}

var _ error = &ErrorResponse{}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("CASP Signing Error: Title: %s, Details: %s, Status: %d", e.Title, e.Details, e.Status)
}
