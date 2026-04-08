package delete_user_function

import "net/http"

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User Deleted"))
}
