/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

func Get(url string, target interface{}) error {
	// a body is going to be closed later
	//nolint:bodyclose
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}

	defer func(body io.Closer) {
		innerErr := body.Close()
		if err != nil {
			panic(innerErr)
		}
	}(r.Body)

	if err = json.NewDecoder(r.Body).Decode(target); err == io.EOF {
		return nil
	}

	return err
}
