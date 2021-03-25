package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckJSONRequestType(t *testing.T) {
	tests := map[string]struct {
		contType string
		code     int
	}{
		"Bad content-type":  {"text/plain", 400},
		"Good content-type": {"application/json", 200},
	}

	srv := &Server{}
	middleWare := srv.CheckJSONRequestType(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			wr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "http://127.0.0.1", nil)
			req.Header.Set("Content-type", tc.contType)

			middleWare(wr, req)

			if wr.Result().StatusCode != tc.code {
				t.Fatalf("Error! Expected %v, got %v", wr.Result().StatusCode, tc.code)
			}

		})
	}
}

func TestJSONHeader(t *testing.T) {
	srv := GetTestServer()

	defer srv.Close()

	res, err := http.Get(fmt.Sprintf("%s/not_found", srv.URL))

	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	if res.Header.Get("Content-type") != "application/json" {
		t.Fatalf("Expected %v, but got %v", "application/json", res.Header.Get("Content-type"))
	}
}
