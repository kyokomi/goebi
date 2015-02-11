package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/kyokomi/goebi/goebi"
	"github.com/kyokomi/goebi/goebi/notice"
)

var ebi *goebi.Client

func main() {

	var hostName, apiPath, apiKey string
	flag.StringVar(&hostName, "host", "http://localhost:3000", "errbit hostName.")
	flag.StringVar(&apiPath, "path", "/api/v3/projects", "errbit apiPath.")
	flag.StringVar(&apiKey, "key", "", "errbit app apikey.")
	flag.Parse()

	ebi = goebi.New(goebi.Options{
		Host:    hostName,
		APIPath: apiPath,
		APIKey:  apiKey,
	})

	serve()
}

func serve() {
	fmt.Println("serve start")

	http.HandleFunc("/", index)

	http.ListenAndServe(":8000", nil)
}

func index(_ http.ResponseWriter, r *http.Request) {

	var n *notice.Notice
	n = goebi.NewNotice(errors.New("Test error"))
	n.SetUserInfo(notice.User{
		UserID:       111111111,
		UserName:     "test name",
		UserUsername: "kyokomi",
		UserEmail:    "example@mail.com",
	})
	n.SetWhere("example", "index")
	n.SetEnvRuntime()
	n.SetProfiles()
	n.SetHTTPRequest(r)

	if err := ebi.SendNotice(n); err != nil {
		log.Fatalln(err)
	}
}
