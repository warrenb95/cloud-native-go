package api

import "net/http"

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello gorilla/mux\n"))
}
