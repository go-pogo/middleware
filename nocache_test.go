// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package middleware

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNoCache(t *testing.T) {
	t.Run("headers present", func(t *testing.T) {
		rec := httptest.NewRecorder()
		NoCache(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		})).ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))

		for k, v := range noCacheHeaders {
			assert.Equalf(t, v, rec.Header().Get(k), "%s header not set by middleware.", k)
		}
	})
	t.Run("headers removed", func(t *testing.T) {

	})
}
