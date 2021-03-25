package main

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/VladimirStepanov/urlshortener/pkg/config"
	"github.com/VladimirStepanov/urlshortener/pkg/shortener/base62"
	"github.com/VladimirStepanov/urlshortener/pkg/store"
	"github.com/VladimirStepanov/urlshortener/pkg/store/teststore"
	"github.com/sirupsen/logrus"
)

var (
	defaultItem = &store.Item{ID: 284772472784, BaseItem: store.BaseItem{URL: "https://vk.com", Visits: 100, Expire: "10.1.2380 1:0:0", Once: false}}

	deleteItem = &store.Item{ID: 431816, BaseItem: store.BaseItem{URL: "https://vk.com", Visits: 100, Expire: "10.1.2380 1:0:0", Once: false}}

	defaultItemWithAlreadyOnce = &store.Item{ID: 25433331007, BaseItem: store.BaseItem{URL: "https://vk.com", Visits: 1, Expire: "10.1.2380 1:0:0", Once: true}}

	defaultResponse = &ResponseItem{"Ubrm0af", store.BaseItem{URL: "https://vk.com", Visits: 100, Expire: "10.1.2380 1:0:0", Once: false}}

	expiredItem = &store.Item{ID: 111111, BaseItem: store.BaseItem{URL: "https://vk.com", Visits: 100, Expire: "10.1.1994 1:0:0", Once: true}}
)

//CheckFatal - check err. if it not nil, call t.Fatal
func CheckFatal(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

//GetTestMap ...
func GetTestMap() map[uint64]*store.Item {
	return map[uint64]*store.Item{
		defaultItem.ID:                defaultItem,
		expiredItem.ID:                expiredItem,
		defaultItemWithAlreadyOnce.ID: defaultItemWithAlreadyOnce,
		deleteItem.ID:                 deleteItem,
	}
}

//GetTestServer ...
func GetTestServer() *httptest.Server {
	log := &logrus.Logger{}
	store := teststore.New(GetTestMap())
	conf := &config.Config{}
	s := &Server{log, store, conf, base62.New()}
	s.log.SetOutput(ioutil.Discard)
	srv := httptest.NewServer(s.router())
	return srv
}
