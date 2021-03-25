package main

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	"github.com/VladimirStepanov/urlshortener/pkg/store"
)

//EncodeResponse ...
type EncodeResponse struct {
	Status string `json:"status"`
	URL    string `json:"url"`
}

//ResponseItem - response json data for GET request
type ResponseItem struct {
	ID string `json:"id"`
	store.BaseItem
}

//Response struct for json response
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (s *Server) serverError(w http.ResponseWriter, err error) {
	s.log.Errorf("Internal error: %v %s", err, string(debug.Stack()))
	http.Error(w, `{"status": "error", "message": "Internal server error"}`, http.StatusInternalServerError)
}

//ResponseJSON ...
func (s *Server) ResponseJSON(w http.ResponseWriter, r interface{}, statusCode int) {
	resp, err := json.Marshal(r)

	if err != nil {
		s.serverError(w, err)
	}

	w.WriteHeader(statusCode)

	w.Write(resp)
}

func (s *Server) response404(w http.ResponseWriter, r *http.Request) {
	s.ResponseJSON(w, &Response{"error", "page not found"}, http.StatusNotFound)
}
