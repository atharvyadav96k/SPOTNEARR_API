package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/atharvyadav96k/spotnearr/search-svc/applayer"
	"github.com/atharvyadav96k/spotnearr/search-svc/config"
)

func main() {
	app := applayer.Init()
	fmt.Printf("Search Service started at http://localhost:%s\n", config.C.Port)
	log.Fatal(http.ListenAndServe(":"+config.C.Port, app.NewMux()))
}
