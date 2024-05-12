// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package middleware

import (
	"net/http"
)

type Wrapper interface {
	// Wrap wraps the [http.Handler] with additional logic.
	Wrap(next http.Handler) http.Handler
}

var _ Wrapper = (*WrapperFunc)(nil)

type WrapperFunc func(next http.Handler) http.Handler

func (fn WrapperFunc) Wrap(next http.Handler) http.Handler { return fn(next) }

// Wrap [http.Handler] h with additional logic via the provided [Wrapper]s.
func Wrap(handler http.Handler, wrap ...Wrapper) http.Handler {
	for i := len(wrap) - 1; i >= 0; i-- {
		handler = wrap[i].Wrap(handler)
	}
	return handler
}

var _ Wrapper = (*Middleware)(nil)

// Middleware consists of multiple [Wrapper]s which can be wrapped around a
// [http.Handler] using [Wrap].
type Middleware []Wrapper

func (m Middleware) Wrap(next http.Handler) http.Handler {
	return Wrap(next, m...)
}

var _ Wrapper = (*HandlerFunc)(nil)

// HandlerFunc is a [http.Handler] which implements the [Wrapper] interface.
type HandlerFunc http.HandlerFunc

func (fn HandlerFunc) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		fn(wri, req)
		next.ServeHTTP(wri, req)
	})
}
