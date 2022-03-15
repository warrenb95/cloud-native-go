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
	store := store.New(make(map[string]interface{}))
	server := api.RESTServer{
		Store: store,
	}

	r.HandleFunc("/", server.IndexHandler)
	r.HandleFunc("/v1/{key}", server.PutKeyValueHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", server.GetKeyValueHandler).Methods("GET")
	r.HandleFunc("/v1/{key}", server.DeleteKeyValueHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
