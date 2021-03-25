package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/VladimirStepanov/urlshortener/pkg/store"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gorilla/mux"
)

//EncodeRequest - POST data
type EncodeRequest struct {
	URL    string `json:"url"`
	Expire string `json:"expire"`
	Once   bool   `json:"once"`
}

//EncodeURL ...
func (s *Server) EncodeURL(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	er := EncodeRequest{}
	err := dec.Decode(&er)

	if err != nil {
		s.ResponseJSON(w, &Response{"error", "bad json"}, 400)
		return
	}

	err = validation.ValidateStruct(&er,
		validation.Field(&er.URL, validation.Required.Error("is required"), is.URL.Error("invalid url")),
		validation.Field(&er.Expire, validation.Required.Error("is required"), validation.Date("2.1.2006 15:4:5").Error("invalid date")),
	)

	if err != nil {
		s.ResponseJSON(w, &Response{"error", err.Error()}, 400)
		return
	}

	dt, _ := time.Parse("2.1.2006 15:4:5", er.Expire)

	id, err := s.db.Save(er.URL, dt, er.Once)

	if err != nil {
		if err == store.ErrExpired {
			s.ResponseJSON(w, &Response{"error", "expire: date is expired."}, 400)
			return
		}
		s.serverError(w, err)
		return
	}

	s.ResponseJSON(w, &EncodeResponse{"success", fmt.Sprintf("http://%s:%s/%s", s.config.Host, s.config.Port, s.shortener.Encode(id))}, 200)
}

//GetInfoHandler ...
func (s *Server) GetInfoHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id, err := s.shortener.Decode(vars["id"])

	if err != nil {
		s.serverError(w, err)
		return
	}

	item, err := s.db.Load(id)

	if err != nil {
		if err == store.ErrItemNotFound {
			s.response404(w, r)
			return
		}
		s.serverError(w, err)
		return
	}

	respItem := ResponseItem{
		ID: vars["id"], BaseItem: store.BaseItem{
			URL: item.URL, Visits: item.Visits, Expire: item.Expire, Once: item.Once,
		},
	}

	s.ResponseJSON(w, respItem, 200)

}

//RedirectURL - redirect to original URL
func (s *Server) RedirectURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := s.shortener.Decode(vars["id"])

	if err != nil {
		s.response404(w, r)
		return
	}

	item, err := s.db.Load(id)

	if err != nil {
		if err == store.ErrItemNotFound {
			s.response404(w, r)
			return
		}
		s.serverError(w, err)
		return
	}

	if item.Once && item.Visits > 0 {
		s.response404(w, r)
		return
	}

	err = s.db.IncVisits(id)
	if err != nil {
		s.response404(w, r)
		return
	}

	http.Redirect(w, r, item.URL, http.StatusFound)
}

//DeleteURL - delete URL from database
func (s *Server) DeleteURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := s.shortener.Decode(vars["id"])

	if err != nil {
		s.response404(w, r)
		return
	}

	_, err = s.db.Remove(id)

	if err != nil {
		if err == store.ErrItemNotFound {
			s.response404(w, r)
			return
		}
		s.serverError(w, err)
		return
	}

	s.ResponseJSON(w, struct {
		Status string `json:"status"`
	}{"success"}, 200)
}
