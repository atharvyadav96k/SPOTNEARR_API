package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/atharvyadav96k/spotnearr/vendor-svc/applayer"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/config"
)

func main() {
	app := applayer.Init()
	fmt.Printf("Vendor Service started at http://localhost:%s\n", config.C.Port)
	log.Fatal(http.ListenAndServe(":"+config.C.Port, app.NewMux()))
}
