package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/warrenb95/cloud-native-go/internal/api"
	"github.com/warrenb95/cloud-native-go/internal/store"
)

func main() {
	r := mux.NewRouter()
	memStore := store.New(make(map[string]string))
	logger, err := initTransactionLogger(memStore)
	if err != nil {
		log.Fatalf("cannot load from transaction logger: %v", err)
	}
	server := api.New(memStore, logger)

	r.HandleFunc("/", server.IndexHandler)
	r.HandleFunc("/v1/{key}", server.PutKeyValueHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", server.GetKeyValueHandler).Methods("GET")
	r.HandleFunc("/v1/{key}", server.DeleteKeyValueHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServeTLS(":8080", "localhost.csr", "localhost.key", r))
}

func initTransactionLogger(memStore *store.Store) (api.TransactionLogger, error) {
	logger, err := store.NewPostgresTransactionLogger(store.PostgresConfig{DBName: "testdb", Host: "localhost", Port: "5432", User: "postgres", Password: "password"})
	if err != nil {
		return nil, fmt.Errorf("failed to create event logger: %w", err)
	}

	events, errors := logger.ReadEvents()
	e, ok := store.Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case store.EventDelete:
				err = memStore.Delete(e.Key)
			case store.EventPut:
				err = memStore.Put(e.Key, string(e.Value))
			}
		}
	}

	logger.Run()

	return logger, err
}
