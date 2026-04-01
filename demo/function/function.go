package demofunction

import "net/http"

func Demo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
