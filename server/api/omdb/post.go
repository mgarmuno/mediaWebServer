package omdb

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/mgarmuno/mediaWebServer/server/items"
)

const (
	url = "https://www.omdbapi.com/?"
)

type OmdbResponse struct {
	TotalResults string
	Response     string
	Search       []items.Movie
}

func doPost(w http.ResponseWriter, r *http.Request) {

	var finished bool = false
	var page int = 1
	var totalMoviesGetted int = 0
	movie := getMovieInfoFromRequest(r)
	var completeMoviesResponded []items.Movie

	for !finished {
		req := prepareRequestQueryByMovie(movie, page)
		page++

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error calling endpoint OMDB", err)
			return
		}
		defer resp.Body.Close()

		movies := getDecodedResponseBody(resp)
		completeMoviesResponded = append(completeMoviesResponded, movies.Search...)
		totalMoviesGetted = totalMoviesGetted + len(movies.Search)
		moviesInResponse, err := strconv.Atoi(movies.TotalResults)
		if err != nil || moviesInResponse <= totalMoviesGetted {
			finished = true
			log.Println("Loop finished")
		}
	}
	fmt.Println("Total movies getted from OMDB", len(completeMoviesResponded))
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	moviesResponse := &OmdbResponse{
		Response:     "True",
		TotalResults: strconv.Itoa(len(completeMoviesResponded)),
		Search:       completeMoviesResponded}
	data, err := json.Marshal(moviesResponse)
	if err != nil {
		data, _ := json.Marshal(&OmdbResponse{Response: "False", TotalResults: "0", Search: nil})
		fmt.Fprint(w, string(data))
		return
	}
	fmt.Fprint(w, string(data))
}

func getDecodedResponseBody(resp *http.Response) OmdbResponse {
	decoder := json.NewDecoder(resp.Body)

	var movies OmdbResponse
	err := decoder.Decode(&movies)
	if err != nil {
		log.Println("Error decoding response:", err)
	}

	return movies
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

func prepareRequestQueryByMovie(movie items.Movie, page int) *http.Request {
	req, err := http.NewRequest("POST", url, nil)

	if err != nil {
		log.Println("Error perparing OMDB request:", err)
		return req
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	q := req.URL.Query()
	q.Add("apikey", apiKey)
	q.Add("s", movie.Title)
	q.Add("page", strconv.Itoa(page))
	if movie.Year != "" {
		q.Add("y", movie.Year)
	}

	req.URL.RawQuery = q.Encode()
	return req
}
