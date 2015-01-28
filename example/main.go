package main

import (
	"github.com/kyokomi/gobit/gobit"
	"flag"
	"log"
	"errors"
)

func main() {
	
	var hostName, apiPath, apiKey string
	flag.StringVar(&hostName, "host", "http://localhost:3000", "errbit hostName.")
	flag.StringVar(&apiPath, "path", "/api/v3/projects", "errbit apiPath.")
	flag.StringVar(&apiKey, "key", "", "errbit app apikey.")
	flag.Parse()
	
	c := gobit.New(gobit.Options{
		Host: hostName,
		ApiPath: apiPath,
		ApiKey: apiKey,
	})
	
	n := gobit.NewNotice(errors.New("errorだよ!!"), nil)
	
	if err := c.Send(n); err != nil {
		log.Fatalln(err)
	}
}
