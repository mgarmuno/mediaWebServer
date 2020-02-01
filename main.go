package main

import (
	"fmt"
	"net/http"

	"github.com/mgarmuno/mediaWebServer/server/data"
)

func main() {
	initialChecks()
	http.HandleFunc("/", indexHandler)

	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
}

func initialChecks() {
	data.OpenConnection()
}
