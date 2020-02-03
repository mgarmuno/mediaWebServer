package omdb

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

const (
	apiKey string = "daee70b3"
)

type OmdbAPI struct{}

var lock sync.Mutex

func (omdb *OmdbAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))

	switch r.Method {
	case http.MethodPost:
		doPost(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unsupported method '%v' to %v\n", r.Method, r.URL)
		log.Printf("Unsupported method '%v' to %v\n", r.Method, r.URL)
	}
}
