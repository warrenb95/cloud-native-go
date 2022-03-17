package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/warrenb95/cloud-native-go/internal/api"
	"github.com/warrenb95/cloud-native-go/internal/persistance"
	"github.com/warrenb95/cloud-native-go/internal/store"
)

func main() {
	r := mux.NewRouter()
	store := store.New(make(map[string]string))
	logger, err := initTransactionLogger(store)
	if err != nil {
		log.Fatalf("cannot load from transaction logger: %v", err)
	}
	server := api.New(store, logger)

	r.HandleFunc("/", server.IndexHandler)
	r.HandleFunc("/v1/{key}", server.PutKeyValueHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", server.GetKeyValueHandler).Methods("GET")
	r.HandleFunc("/v1/{key}", server.DeleteKeyValueHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func initTransactionLogger(store *store.Store) (*persistance.FileTransactionLogger, error) {
	logger, err := persistance.NewFileTransactionLogger("transaction.log")
	if err != nil {
		return nil, fmt.Errorf("failed to create event logger: %w", err)
	}

	events, errors := logger.ReadEvents()
	e, ok := persistance.Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case persistance.EventDelete:
				err = store.Delete(e.Key)
			case persistance.EventPut:
				err = store.Put(e.Key, string(e.Value))
			}
		}
	}

	logger.Run()

	return logger, err
}
