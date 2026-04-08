package users_function

import "net/http"

func CreateUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Users"))
}
