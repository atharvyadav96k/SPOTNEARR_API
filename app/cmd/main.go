package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/applayer"
	"github.com/atharvyadav96k/SPOTNEARR_API/config"
)

func main() {
	app := applayer.Init()
	fmt.Printf("Server started at http://localhost:%s\n", config.C.Port)
	log.Fatal(http.ListenAndServe(":"+config.C.Port, app.NewMux()))
}
