package main

import (
	"fmt"

	"github.com/VladimirStepanov/urlshortener/pkg/config"
	"github.com/VladimirStepanov/urlshortener/pkg/shortener/base62"
	"github.com/VladimirStepanov/urlshortener/pkg/store/redis"
)

func main() {
	conf, err := config.New(".env")

	if err != nil {
		fmt.Println("Error while create conf instance", err)
	}

	serv, err := New(conf, redis.New(conf), base62.New())

	if err != nil {
		fmt.Println("Error while create Server instance", err)
	}

	if err = serv.Start(); err != nil {
		fmt.Println(err)
	}

}
