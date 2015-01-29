package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/guregu/kami"
	"github.com/kyokomi/gobit/gobit"
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
		Host:    hostName,
		APIPath: apiPath,
		APIKey:  apiKey,
	})

	serve()
}

func serve() {
	fmt.Println("serve start")

	kami.Get("/", func(_ context.Context, _ http.ResponseWriter, r *http.Request) {

		n := gobit.NewNotice(errors.New("errorだよ"))
		n.SetEnvRuntime()
		n.SetHTTPRequest(r)

		if err := errbit.SendNotice(n); err != nil {
			log.Fatalln(err)
		}
	})

	kami.Serve()
}
