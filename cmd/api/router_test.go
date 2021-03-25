package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestEncodeURLHandler(t *testing.T) {
	tests := map[string]struct {
		data string
		code int
		resp *Response
	}{
		"Success create":                  {`{"url": "https://vk.com", "expire": "10.1.2380 1:0:0"}`, 200, nil},
		"Date is expired":                 {`{"url": "https://vk.com", "expire": "10.1.1984 1:0:0"}`, 400, &Response{"error", "expire: date is expired."}},
		"Url is required":                 {`{"expire": "10.1.2380 1:0:0"}`, 400, &Response{"error", "url: is required."}},
		"Invalid url":                     {`{"url": "bad_url", "expire": "10.1.2380 1:0:0"}`, 400, &Response{"error", "url: invalid url."}},
		"Invalid date":                    {`{"url": "https://vk.com", "expire": "10.1 1:0:0"}`, 400, &Response{"error", "expire: invalid date."}},
		"Expire is required":              {`{"url": "https://vk.com"}`, 400, &Response{"error", "expire: is required."}},
		"Bad json request[bad url type]":  {`{"url": 123, "expire": "10.1.2380 1:0:0"}`, 400, &Response{"error", "bad json"}},
		"Bad json request[unknown field]": {`{"hello": "world"}`, 400, &Response{"error", "bad json"}},
	}

	srv := GetTestServer()
	defer srv.Close()

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest("POST", fmt.Sprintf("%s/encode", srv.URL), bytes.NewReader([]byte(tc.data)))
			CheckFatal(t, err)

			req.Header.Set("Content-type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			CheckFatal(t, err)
			defer resp.Body.Close()

			if resp.StatusCode != tc.code {
				t.Fatalf("Error! Expected code %v, got %v", tc.code, resp.StatusCode)
			}

			if resp.StatusCode != 200 {
				data, err := ioutil.ReadAll(resp.Body)
				CheckFatal(t, err)

				r := Response{}
				err = json.Unmarshal(data, &r)
				CheckFatal(t, err)
				if !reflect.DeepEqual(tc.resp, &r) {
					t.Fatalf("Error! Expected response %v, got %v", tc.resp, r)
				}
			}

		})
	}
}

func TestGetInfoHandler(t *testing.T) {
	tests := map[string]struct {
		url  string
		code int
		item *ResponseItem
	}{
		"Item is found":          {"info/Ubrm0af", http.StatusOK, defaultResponse},
		"Item not found":         {"info/notFound", http.StatusNotFound, nil},
		"Expired item not found": {"info/h4C", http.StatusNotFound, nil},
	}

	srv := GetTestServer()
	defer srv.Close()

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := http.Get(fmt.Sprintf("%s/%s", srv.URL, tc.url))
			CheckFatal(t, err)

			defer res.Body.Close()

			data, err := ioutil.ReadAll(res.Body)
			CheckFatal(t, err)

			if tc.code == http.StatusOK && tc.code == res.StatusCode {
				respItem := ResponseItem{}
				err = json.Unmarshal(data, &respItem)
				CheckFatal(t, err)

				if !reflect.DeepEqual(tc.item, &respItem) {
					t.Fatalf("Error! Expected: %v, got %v\n", tc.item, respItem)
				}
			} else if tc.code != res.StatusCode {
				t.Fatalf("Error! Expected code %v, got %v", tc.code, res.StatusCode)
			}
		})
	}
}

func TestRedirectURLHandler(t *testing.T) {
	srv := GetTestServer()
	defer srv.Close()

	tests := map[string]struct {
		encodedURL string
		code       int
	}{
		"Success redirect":        {"Ubrm0af", http.StatusFound},
		"URL not found":           {"Ub", http.StatusNotFound},
		"Once is already visited": {"poPnVB", http.StatusNotFound},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}}

			resp, err := client.Get(fmt.Sprintf("%s/%s", srv.URL, tc.encodedURL))
			CheckFatal(t, err)

			defer resp.Body.Close()

			if resp.StatusCode != tc.code {
				t.Fatalf("Error! Expected code %v, got %v", tc.code, resp.StatusCode)
			}

		})
	}
}

func TestDeleteURLHandler(t *testing.T) {
	srv := GetTestServer()
	defer srv.Close()

	tests := map[string]struct {
		encodedURL string
		code       int
	}{
		"URL not found":  {"Ub", http.StatusNotFound},
		"Success delete": {"WuYb", http.StatusOK},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := http.Get(fmt.Sprintf("%s/%s", srv.URL, tc.encodedURL))
			CheckFatal(t, err)

			defer resp.Body.Close()

			if resp.StatusCode != tc.code {
				t.Fatalf("Error! Expected code %v, got %v", tc.code, resp.StatusCode)
			}

		})
	}
}
