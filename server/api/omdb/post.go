package omdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mgarmuno/mediaWebServer/server/items"
)

const (
	url = "https://www.omdbapi.com/?"
)

func doPost(w http.ResponseWriter, r *http.Request) {

	movie := getMovieInfoFromRequest(r)
	req := prepareRequestQueryByMovie(movie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error calling endpoint OMDB", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response OMDB API status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprint(w, string(body))
}

func getMovieInfoFromRequest(r *http.Request) items.Movie {
	decoder := json.NewDecoder(r.Body)

	var movie items.Movie
	err := decoder.Decode(&movie)
	if err != nil {
		log.Println("Error decoding response:", err)
	}

	return movie
}

func prepareRequestQueryByMovie(movie items.Movie) *http.Request {
	req, err := http.NewRequest("POST", url, nil)

	if err != nil {
		log.Println("Error perparing OMDB request:", err)
		return req
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	q := req.URL.Query()
	q.Add("apikey", apiKey)
	q.Add("s", movie.Title)
	if movie.Year != "" {
		q.Add("y", movie.Year)
	}

	req.URL.RawQuery = q.Encode()
	return req
}
