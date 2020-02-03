package file

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type FileUploadAPI struct{}

var lock sync.Mutex

func (f FileUploadAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))

	switch r.Method {
	case http.MethodPost:
		doPost(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unsupported methos '%v' to %v\n", r.Method, r.URL)
		log.Printf("Unsupported methos '%v' to %v\n", r.Method, r.URL)
	}
}
