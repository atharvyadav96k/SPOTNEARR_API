package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/atharvyadav96k/SPOTNEARR_API/applayer"
)

func main() {
	app := applayer.Init()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server started at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, app.NewMux()))
}
