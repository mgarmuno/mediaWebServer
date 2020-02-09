package file

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/mgarmuno/mediaWebServer/server/api/omdb"
	"github.com/mgarmuno/mediaWebServer/server/data"
	"github.com/mgarmuno/mediaWebServer/server/items"
	"github.com/mgarmuno/mediaWebServer/server/utils"
)

const (
	url = "http://localhost:8080/api/omdb/"
)

type FileUploadResponse struct {
	Saved    bool
	Filename string
}

func doPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "multipart/form-data")
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("Error uploading file:", err)
		return
	}
	defer file.Close()
	errFw := writeFileOnDisk(handler, file)
	if errFw != nil {
		log.Println("Error saving file on disk:", err)
		return
	}
	movieName, session, episode := utils.GetMovieNameSessionEpisodeByFileName(handler.Filename)
	movies := getCandidateMovies(movieName)
	checkResponse(&w, movies, handler.Filename, session, episode)
	respondFileSaved(w, handler)
}

func writeFileOnDisk(handler *multipart.FileHeader, file multipart.File) error {
	osFile, err := os.Create("/tmp/" + handler.Filename)
	if err != nil {
		return err
	}
	defer osFile.Close()
	hddFile, err := os.OpenFile(osFile.Name(), os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer hddFile.Close()
	_, err = io.Copy(hddFile, file)
	if err != nil {
		return err
	}
	return nil
}

func respondFileSaved(w http.ResponseWriter, handler *multipart.FileHeader) {
	fmt.Println("File saved:", handler.Filename)
	w.WriteHeader(200)
	data := getMovieSavesResponse(handler.Filename)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(data))
}

func getMovieSavesResponse(filename string) []byte {
	uploadResponse := &FileUploadResponse{Saved: true, Filename: filename}
	data, err := json.Marshal(uploadResponse)
	if err != nil {
		log.Println("Error parsing movie response to JSON:", err)
	}
	return data
}

func getCandidateMovies(movieName string) omdb.OmdbResponse {
	fmt.Println("Searching for movie:", movieName)
	req := getRequestQuery(movieName, "")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error getting cadidate movies:", err)
	}
	defer response.Body.Close()
	return omdb.GetDecodedResponseBody(response)
}

func getRequestQuery(movieName string, imdbID string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error consuming API OMDB from API FILE:", err)
		return nil
	}
	q := req.URL.Query()
	if movieName != "" {
		q.Add("title", movieName)
	} else if imdbID != "" {
		q.Add("imdbid", imdbID)
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")
	return req
}

func checkResponse(w *http.ResponseWriter, movies omdb.OmdbResponse, filename string, session string, episode string) {
	if movies.Response == "True" {
		if movies.TotalResults == "1" {
			uploadMovies(movies.Search[0])
		} else {
			// TODO not by ID
		}
	} else {
		// TODO fail in the response
	}
}

func uploadMovies(movie items.Movie) {
	req := getRequestQuery("", movie.ImdbID)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error getting cadidate movies:", err)
	}
	defer response.Body.Close()
	writeFinalFileOnDisk()
	data.InsertMovie(movie)
}

func writeFinalFileOnDisk() {

}
