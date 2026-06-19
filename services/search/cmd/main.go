package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/atharvyadav96k/spotnearr/search-svc/applayer"
	"github.com/atharvyadav96k/spotnearr/search-svc/config"
)

func main() {
	app := applayer.Init()
	srv := &http.Server{
		Addr:              ":" + config.C.Port,
		Handler:           app.NewMux(),
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       4 * time.Second,
		WriteTimeout:      4 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	fmt.Printf("Search Service started at http://localhost:%s\n", config.C.Port)
	log.Fatal(srv.ListenAndServe())
}
