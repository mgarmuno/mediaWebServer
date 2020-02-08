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

func doGet(w http.ResponseWriter, r *http.Request) {
	finish, page, totalMovSum, movie, compMovColl := preparePostVariables(r)

	for !finish {
		compMovColl = doPostEncapsulated(&movie, &page, compMovColl, &totalMovSum, &finish)
		if compMovColl == nil {
			finish = true
			setFailResponse(w)
			return
		}
	}

	prepareSuccessResponse(w, compMovColl)
}

func preparePostVariables(r *http.Request) (bool, int, int, items.Movie, []items.Movie) {
	var finish bool = false
	var page int = 1
	var totalMovSum int = 0
	movie := getMovieInfoFromRequest(r)
	var compMovColl []items.Movie

	return finish, page, totalMovSum, movie, compMovColl
}

func doPostEncapsulated(movie *items.Movie, page *int, completeMoviesResponded []items.Movie, totalMoviesGetted *int, finished *bool) []items.Movie {
	req := getRequestQueryByMovie(movie, page)
	*page++
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	movies := getDecodedResponseBody(resp)
	completeMoviesResponded = append(completeMoviesResponded, movies.Search...)
	*totalMoviesGetted = *totalMoviesGetted + len(movies.Search)
	moviesInResponse, err := strconv.Atoi(movies.TotalResults)

	if err != nil || moviesInResponse <= *totalMoviesGetted {
		*finished = true
		log.Println("Loop finished")
	}

	return completeMoviesResponded
}

func prepareSuccessResponse(w http.ResponseWriter, compMovColl []items.Movie) {
	fmt.Println("Total movies getted from OMDB", len(compMovColl))
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")

	moviesResponse := &OmdbResponse{
		Response:     "True",
		TotalResults: strconv.Itoa(len(compMovColl)),
		Search:       compMovColl}
	data, err := json.Marshal(moviesResponse)

	if err != nil {
		setFailResponse(w)
		return
	}

	fmt.Fprint(w, string(data))
}

func setFailResponse(w http.ResponseWriter) {
	fmt.Println("Somthing went wrong getting the movies search...")
	w.WriteHeader(200)
	data, _ := json.Marshal(&OmdbResponse{Response: "False", TotalResults: "0", Search: nil})

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

func getRequestQueryByMovie(movie *items.Movie, page *int) *http.Request {
	req, err := http.NewRequest("POST", url, nil)

	if err != nil {
		log.Println("Error perparing OMDB request:", err)
		return req
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	q := req.URL.Query()
	q.Add("apikey", apiKey)
	q.Add("s", movie.Title)
	q.Add("page", strconv.Itoa(*page))

	if movie.Year != "" {
		q.Add("y", movie.Year)
	}
	req.URL.RawQuery = q.Encode()

	return req
}
