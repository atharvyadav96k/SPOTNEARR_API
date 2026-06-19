package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/atharvyadav96k/spotnearr/vendor-svc/applayer"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/config"
)

func main() {
	app := applayer.Init()
	srv := &http.Server{
		Addr:              ":" + config.C.Port,
		Handler:           app.NewMux(),
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	fmt.Printf("Vendor Service started at http://localhost:%s\n", config.C.Port)
	log.Fatal(srv.ListenAndServe())
}
