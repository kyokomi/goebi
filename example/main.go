package main

import (
	"github.com/kyokomi/gobit/gobit"
	"flag"
	"log"
	"errors"
	"fmt"
	"net/http"
	"github.com/guregu/kami"
	"golang.org/x/net/context"
)


var errbit *gobit.Client

func main() {
	
	var hostName, apiPath, apiKey string
	flag.StringVar(&hostName, "host", "http://localhost:3000", "errbit hostName.")
	flag.StringVar(&apiPath, "path", "/api/v3/projects", "errbit apiPath.")
	flag.StringVar(&apiKey, "key", "", "errbit app apikey.")
	flag.Parse()

	errbit = gobit.New(gobit.Options{
		Host: hostName,
		ApiPath: apiPath,
		ApiKey: apiKey,
	})
	
	serve()
}

func serve() {
	fmt.Println("serve start")

	kami.Get("/", func(_ context.Context, _ http.ResponseWriter, r *http.Request) {
		
		n := gobit.NewNotice(errors.New("errorだよ"), r)
		if err := errbit.Send(n); err != nil {
			log.Fatalln(err)
		}
	})
	
	kami.Serve()
}


