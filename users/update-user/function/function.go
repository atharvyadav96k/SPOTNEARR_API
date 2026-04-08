package update_user_function

import "net/http"

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User Updated"))
}
