package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/warrenb95/cloud-native-go/internal/api"
	"github.com/warrenb95/cloud-native-go/internal/store"
)

func main() {
	r := mux.NewRouter()
	store := make(store.Store)
	server := api.RESTServer{
		Store: store,
	}

	r.HandleFunc("/", server.IndexHandler)
	r.HandleFunc("/v1/{key}", server.PUTKeyValueHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", server.GetKeyValueHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
