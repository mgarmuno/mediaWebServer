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
	IsByID       bool
}

func doGet(w http.ResponseWriter, r *http.Request) {
	finish, page, totalMovSum, movie, compMovColl, isByID := preparePostVariables(r)
	fmt.Println("Searching movie:", movie)
	for !finish && !isByID {
		compMovColl, isByID = doPostEncapsulated(&movie, &page, compMovColl, &totalMovSum, &finish)
		if compMovColl == nil {
			finish = true
			setFailResponse(w, nil)
			return
		}
	}
	prepareSuccessResponse(w, compMovColl, isByID)
}

func preparePostVariables(r *http.Request) (bool, int, int, items.Movie, []items.Movie, bool) {
	var finish bool = false
	var page int = 1
	var totalMovSum int = 0
	movie := getMovieInfoFromRequest(r)
	var compMovColl []items.Movie
	return finish, page, totalMovSum, movie, compMovColl, false
}

func doPostEncapsulated(movie *items.Movie, page *int, completeMoviesResponded []items.Movie, totalMoviesGetted *int, finished *bool) ([]items.Movie, bool) {
	req, isByID := getRequestQueryByMovie(movie, page)
	*page++
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, isByID
	}
	defer resp.Body.Close()
	var movies OmdbResponse
	if isByID {
		movies = getDecodedResponseBodyByID(resp)
	} else {
		movies = GetDecodedResponseBody(resp)
	}
	completeMoviesResponded = append(completeMoviesResponded, movies.Search...)
	*totalMoviesGetted = *totalMoviesGetted + len(movies.Search)
	moviesInResponse, err := strconv.Atoi(movies.TotalResults)
	if err != nil || moviesInResponse <= *totalMoviesGetted {
		*finished = true
		log.Println("Loop finished")
	}
	return completeMoviesResponded, isByID
}

func prepareSuccessResponse(w http.ResponseWriter, compMovColl []items.Movie, isByID bool) {
	fmt.Println("Total movies getted from OMDB", len(compMovColl))
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	moviesResponse := &OmdbResponse{
		Response:     "True",
		TotalResults: strconv.Itoa(len(compMovColl)),
		Search:       compMovColl,
		IsByID:       isByID}
	data, err := json.Marshal(moviesResponse)
	if err != nil {
		setFailResponse(w, err)
		return
	}
	fmt.Fprint(w, string(data))
}

func setFailResponse(w http.ResponseWriter, err error) {
	fmt.Println("Something went wrong getting the movies search:", err)
	w.WriteHeader(200)
	data, _ := json.Marshal(&OmdbResponse{Response: "False", TotalResults: "0", Search: nil})
	fmt.Fprint(w, string(data))
}

func getDecodedResponseBodyByID(res *http.Response) OmdbResponse {
	decoder := json.NewDecoder(res.Body)
	var movie items.Movie
	err := decoder.Decode(&movie)
	if err != nil {
		log.Println("Error decoding response ByID:", err)
	}
	movies := OmdbResponse{
		Response:     "True",
		TotalResults: "1",
		Search:       []items.Movie{movie},
		IsByID:       true}
	return movies
}

func GetDecodedResponseBody(res *http.Response) OmdbResponse {
	decoder := json.NewDecoder(res.Body)
	var movies OmdbResponse
	err := decoder.Decode(&movies)
	if err != nil {
		log.Println("Error decoding response:", err)
	}
	return movies
}

func getMovieInfoFromRequest(r *http.Request) items.Movie {
	var title string = r.URL.Query().Get("title")
	var year string = r.URL.Query().Get("year")
	var imdbID string = r.URL.Query().Get("imdbid")
	movie := items.Movie{
		Title:  title,
		Year:   year,
		ImdbID: imdbID}
	return movie
}

func getRequestQueryByMovie(movie *items.Movie, page *int) (*http.Request, bool) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error perparing OMDB request:", err)
		return req, false
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	var isByID bool = setQueryParameters(req, movie, page)
	return req, isByID
}

func setQueryParameters(req *http.Request, movie *items.Movie, page *int) bool {
	var isByID bool = false
	q := req.URL.Query()
	q.Add("apikey", apiKey)
	if movie.ImdbID != "" {
		q.Add("i", movie.ImdbID)
		q.Add("plot", "full")
		isByID = true
	} else {
		if movie.Title != "" {
			q.Add("s", movie.Title)
		}
		if movie.Year != "" {
			q.Add("y", movie.Year)
		}
		q.Add("page", strconv.Itoa(*page))
	}
	req.URL.RawQuery = q.Encode()
	return isByID
}
