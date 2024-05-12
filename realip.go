// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Portions of this file are cloned from https://github.com/zenazn/goji.
// The original code is licensed under the MIT license and is owned by its
// author(s).

package middleware

import (
	"net/http"
	"strings"
)

var xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
var xRealIP = http.CanonicalHeaderKey("X-Real-IP")

// WithRealIP adds middleware that sets the [RemoteAddr] field of the
// [http.Request] to the value of the client's real IP address, using [RealIP].
func WithRealIP() Wrapper { return WrapperFunc(RealIP) }

// RealIP is a middleware that sets a http.Request's RemoteAddr to the results
// of parsing either the X-Forwarded-For header or the X-Real-IP header (in
// that order).
//
// This middleware should be inserted fairly early in the middleware stack to
// ensure that subsequent layers (e.g., request loggers) which examine the
// RemoteAddr will see the intended value.
//
// You should only use this middleware if you can trust the headers passed to
// you (in particular, the two headers this middleware uses), for example
// because you have placed a reverse proxy like HAProxy or nginx in front of the
// server. If your reverse proxies are configured to pass along arbitrary header
// values from the client, or if you use this middleware without a reverse
// proxy, malicious clients will be able to make you very sad (or, depending on
// how you're using RemoteAddr, vulnerable to an attack of some sort).
//
// This function is based on the RealIP middleware from the Goji web framework.
func RealIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rip := realIP(r); rip != "" {
			r.RemoteAddr = rip
		}
		next.ServeHTTP(w, r)
	})
}

func realIP(r *http.Request) string {
	var ip string

	if xff := r.Header.Get(xForwardedFor); xff != "" {
		i := strings.Index(xff, ", ")
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	} else if xrip := r.Header.Get(xRealIP); xrip != "" {
		ip = xrip
	}

	return ip
}
