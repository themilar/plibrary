package main

import "net/http"

func (app *application) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.Info("logging request", "remote address", r.RemoteAddr, "protocol", r.Proto, "method", r.Method, "uri", r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
