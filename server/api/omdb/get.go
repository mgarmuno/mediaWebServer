package omdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/mgarmuno/mediaWebServer/server/items"
	"github.com/mgarmuno/mediaWebServer/server/utils"
)

const (
	url      = "https://www.omdbapi.com/?"
	notFound = "Movie not found!"
)

type OmdbResponse struct {
	TotalResults string
	Response     string
	Error        string
	Search       []items.Movie
	IsByID       bool
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

func DecodeOmdbResponse(res *http.Response) OmdbResponse {
	decoder, _ := ioutil.ReadAll(res.Body)
	var omdbResponse OmdbResponse
	json.Unmarshal([]byte(decoder), &omdbResponse)
	return omdbResponse
}

func doGet(w http.ResponseWriter, r *http.Request) {
	finish, page, totalMovSum, movie, compMovColl, isByID, titleForSearch := preparePostVariables(r)
	var tries int = 30
	for !finish && tries > 0 {
		fmt.Println("Try number:", tries, "searching for", titleForSearch)
		compMovColl, isByID = doPostEncapsulated(&movie, &page, compMovColl, &totalMovSum, &finish, titleForSearch)
		if titleForSearch == "" {
			finish = true
			setFailResponse(w, nil)
			return
		} else if len(compMovColl) == 0 {
			titleForSearch = utils.RemoveLastWord(titleForSearch)
			if titleForSearch == "" {
				break
			}
		}
		if finish {
			break
		}
		tries--
	}
	prepareSuccessResponse(w, compMovColl, isByID)
}

func preparePostVariables(r *http.Request) (bool, int, int, items.Movie, []items.Movie, bool, string) {
	var finish bool = false
	var page int = 1
	var totalMovSum int = 0
	movie := getMovieInfoFromRequest(r)
	var compMovColl []items.Movie = []items.Movie{}
	return finish, page, totalMovSum, movie, compMovColl, false, movie.Title
}

func doPostEncapsulated(movie *items.Movie, page *int, completeMoviesResponded []items.Movie, totalMoviesGetted *int, finished *bool, titleForSearch string) ([]items.Movie, bool) {
	req, isByID := getRequestQueryByMovie(movie, page, titleForSearch)
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
	if movies.Error == notFound {
		return completeMoviesResponded, isByID
	}
	completeMoviesResponded = append(completeMoviesResponded, movies.Search...)
	*totalMoviesGetted = *totalMoviesGetted + len(movies.Search)
	moviesInResponse, err := strconv.Atoi(movies.TotalResults)
	if err != nil || moviesInResponse <= *totalMoviesGetted || isByID {
		*finished = true
		log.Println("Loop finished")
	}
	*page++
	return completeMoviesResponded, isByID
}

func prepareSuccessResponse(w http.ResponseWriter, compMovColl []items.Movie, isByID bool) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	moviesResponse := OmdbResponse{
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

func getMovieInfoFromRequest(r *http.Request) items.Movie {
	var title string = r.URL.Query().Get("title")
	var year string = r.URL.Query().Get("year")
	var imdbID string = r.URL.Query().Get("imdbid")
	movie := items.Movie{
		Title:  title,
		Year:   year,
		ImdbID: imdbID}
	fmt.Println(movie)
	return movie
}

func getRequestQueryByMovie(movie *items.Movie, page *int, titleForSearch string) (*http.Request, bool) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error perparing OMDB request:", err)
		return req, false
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	var isByID bool = setQueryParameters(req, movie, page, titleForSearch)
	return req, isByID
}

func setQueryParameters(req *http.Request, movie *items.Movie, page *int, titleForSearch string) bool {
	var isByID bool = false
	q := req.URL.Query()
	q.Add("apikey", apiKey)
	if movie.ImdbID != "" {
		q.Add("i", movie.ImdbID)
		q.Add("plot", "full")
		isByID = true
	} else {
		if movie.Title != "" {
			q.Add("s", titleForSearch)
		}
		if movie.Year != "" {
			q.Add("y", movie.Year)
		}
		q.Add("page", strconv.Itoa(*page))
	}
	req.URL.RawQuery = q.Encode()
	return isByID
}
