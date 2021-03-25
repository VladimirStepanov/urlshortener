package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) router() http.Handler {
	mux := mux.NewRouter()

	mux.HandleFunc("/info/{id}", s.GetInfoHandler).Methods("GET")
	mux.HandleFunc("/encode", s.CheckJSONRequestType(s.EncodeURL)).Methods("POST")
	mux.HandleFunc("/{id}", s.RedirectURL).Methods("GET")
	mux.HandleFunc("/{id}", s.DeleteURL).Methods("DELETE")

	mux.NotFoundHandler = http.HandlerFunc(s.response404)
	return s.JSONHeader(s.Log(mux))
}
