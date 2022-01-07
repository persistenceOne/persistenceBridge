/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package rest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

const HTTPTimeout = 10 * time.Second

// nolint fixme move to a call context or a new type
// nolint: gochecknoglobals
var httpClient = &http.Client{Timeout: HTTPTimeout}

func Get(url string, target interface{}) error {
	// nolint a body is going to be closed later
	// nolint:bodyclose
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

	if err = json.NewDecoder(r.Body).Decode(target); errors.Is(err, io.EOF) {
		return nil
	}

	return err
}
