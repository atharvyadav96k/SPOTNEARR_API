package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/atharvyadav96k/SPOTNEARR_API/applayer"
	"github.com/atharvyadav96k/SPOTNEARR_API/config"
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
	fmt.Printf("Server started at http://localhost:%s\n", config.C.Port)
	log.Fatal(srv.ListenAndServe())
}
