package api

import (
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/warrenb95/cloud-native-go/internal/store"
)

type Store interface {
	Put(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}

type RESTServer struct {
	Store Store
}

func (s *RESTServer) IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello gorilla/mux\n"))
}

// PUTKeyValueHandler expects path "/v1/{key}" and will then save that to the store.
func (s *RESTServer) PUTKeyValueHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	err = s.Store.Put(key, string(value))
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetKeyValueHandler expects path "/v1/{key}" and will get the value for the provided key if it exists in the store.
func (s *RESTServer) GetKeyValueHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := s.Store.Get(key)
	if err != nil {
		if errors.Is(err, store.ErrNoSuchKey) {
			http.Error(w,
				err.Error(),
				http.StatusNotFound)
			return
		}
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.Write([]byte(value))
}

// DeleteKeyValueHandler expects path "v1/{key}" and will delete the key value pair from the store.
func (s *RESTServer) DeleteKeyValueHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	err := s.Store.Delete(key)
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}
}
