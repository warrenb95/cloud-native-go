package api

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
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

// expects path "/v1/keyvalue/{key}"
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
