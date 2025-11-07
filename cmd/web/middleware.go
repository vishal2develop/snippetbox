package main

import (
	"fmt"
	"net/http"
)

// middleware to add common headers

func commonHeaders(next http.Handler) http.Handler {
	// create an anonymous function to handle the request
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// any code here will execute on the way down the chain
		// set common headers
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		w.Header().Set("Server", "Go")
		// call the next handler in the chain
		next.ServeHTTP(w, r)

		// any code here will execute on the way up the chain
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	// anonymous function to handle the request
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip        = r.RemoteAddr
			proto     = r.Proto
			method    = r.Method
			uri       = r.URL.RequestURI()
			host      = r.Host
			userAgent = r.UserAgent()
		)

		app.logger.Info("Received request", "ip", ip, "proto", proto, "method", method, "uri", uri, "host", host, "user_agent", userAgent)

		// call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// mioddleware to recover from panics
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event
		// of a panic).
		defer func() {
			pv := recover()
			if pv != nil {
				// set connection close header
				w.Header().Set("Connection", "close")
				// call the serverError() helper to handle the error
				app.serverError(w, r, fmt.Errorf("%v", pv))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
