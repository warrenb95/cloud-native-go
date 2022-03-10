package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/warrenb95/cloud-native-go/internal/store"
	"github.com/warrenb95/cloud-native-go/internal/store/api"
)

func main() {
	r := mux.NewRouter()
	store := make(store.Store)
	server := api.RESTServer{
		Store: store,
	}

	r.HandleFunc("/", server.IndexHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}
