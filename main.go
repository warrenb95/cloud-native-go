package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/warrenb95/cloud-native-go/internal/store/api"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", api.IndexHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}
