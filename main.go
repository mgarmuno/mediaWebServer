package main

import (
	"fmt"
	"net/http"

	"github.com/mgarmuno/mediaWebServer/server/api/movie"
	"github.com/mgarmuno/mediaWebServer/server/api/omdb"
	"github.com/mgarmuno/mediaWebServer/server/data"
)

func main() {
	initialChecks()

	fs := http.FileServer(http.Dir("client"))
	http.Handle("/", fs)
	http.Handle("/api/omdb/", &omdb.OmdbAPI{})
	http.Handle("/api/movie/", &movie.MovieAPI{})

	fmt.Println("Serving...")
	http.ListenAndServe(":8080", nil)
}

func initialChecks() {
	data.CheckDatabase()
}
