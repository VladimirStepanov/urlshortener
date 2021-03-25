package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestGetLogger(t *testing.T) {
	tests := map[string]struct {
		logLevel string
		isError  bool
		err      string
	}{
		"Success":                {"INFO", false, ""},
		"Error! Not valid level": {"LLLLLL", true, `not a valid logrus Level: "LLLLLL"`},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := getLogger(tc.logLevel)
			if tc.isError && err == nil {
				t.Fatalf("Expected error: %v, but got nil", tc.err)
			}

			if tc.isError && err != nil {
				if tc.err != err.Error() {
					t.Fatalf("Expected errror: %v, but got: %v", tc.err, err)
				}
			}
		})
	}

}

func TestNotFoundJSON(t *testing.T) {
	srv := GetTestServer()

	defer srv.Close()

	res, err := http.Get(fmt.Sprintf("%s/not_found", srv.URL))

	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected %v, but got %v", http.StatusNotFound, res.StatusCode)
	}
}
