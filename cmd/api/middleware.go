package main

import (
	"net/http"

	"github.com/VladimirStepanov/urlshortener/pkg/middleware"
)

//CheckJSONRequestType ...
func (s *Server) CheckJSONRequestType(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-type") != "application/json" {
			s.ResponseJSON(w, &Response{"error", `content-type must be "application/json"`}, 400)
			return
		}
		next(w, r)
	}
}

//JSONHeader ...
func (s *Server) JSONHeader(mux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		mux.ServeHTTP(w, r)
	})
}

//Log ...
func (s *Server) Log(mux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lwr := &middleware.LoggerWR{W: w, StatusCode: 200}

		mux.ServeHTTP(lwr, r)
		s.log.Printf("%s request from %s to %s [status code %d]\n", r.Method, r.RemoteAddr, r.RequestURI, lwr.StatusCode)
	})
}
