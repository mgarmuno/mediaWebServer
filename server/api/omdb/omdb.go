package omdb

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type OmdbAPI struct{}
type Omdb struct {
	Title string `json:"title,omitempty"`
	Year  string `json:"year,omitempty"`
}

var db = []*Omdb{}
var lock sync.Mutex

func (omdb *OmdbAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))

	switch r.Method {
	case http.MethodPost:
		doPost(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unsuported method '%v' to %v\n", r.Method, r.URL)
		log.Printf("Unsuported method '%v' to %v\n", r.Method, r.URL)
	}
}
