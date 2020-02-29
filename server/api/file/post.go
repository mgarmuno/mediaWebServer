package file

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/mgarmuno/mediaWebServer/server/api/omdb"
	"github.com/mgarmuno/mediaWebServer/server/data"
	"github.com/mgarmuno/mediaWebServer/server/items"
	"github.com/mgarmuno/mediaWebServer/server/utils"
)

const (
	url        = "http://localhost:8080/api/omdb/"
	tmp        = "/tmp/"
	moviesPath = "/mnt/data/mediaWebServerFiles/movies/"
	imagesPath = "/mnt/data/mediaWebServerFiles/images/"
	sessionAb  = "S"
	episodeAb  = "E"
)

type FileUploadResponse struct {
	Saved    bool
	Filename string
}

type SeveralOptionsResponse struct {
	Saved    bool
	Filename string
	Options  omdb.OmdbResponse
	Episode  string
	Session  string
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
	// respondFileSaved(w, handler)
}

func writeFileOnDisk(handler *multipart.FileHeader, file multipart.File) error {
	osFile, err := os.Create(tmp + handler.Filename)
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
	log.Println("File saved:", handler.Filename)
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
	req := getRequestQuery(movieName, "")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error getting cadidate movies:", err)
	}
	defer response.Body.Close()
	return omdb.GetDecodedResponseBody(response)
}

func getRequestQuery(movieName, imdbID string) *http.Request {
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

func checkResponse(w *http.ResponseWriter, movies omdb.OmdbResponse, filename, session, episode string) {
	if movies.Response == true {
		if movies.TotalResults == "1" {
			uploadMovie(movies.Search[0], filename, session, episode)
		} else {
			respondSeveralMoviesOptions(*w, movies, filename, session, episode)
			fmt.Println("checkResponse", movies)
		}
	} else {
		// TODO fail in the response
	}
}

func uploadMovie(movie items.Movie, filename, session, episode string) {
	req := getRequestQuery("", movie.ImdbID)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error getting cadidate movies:", err)
	}
	defer response.Body.Close()
	omdbResponse := omdb.DecodeOmdbResponse(response)
	finalFilePath := getFinalPathForFile(&movie, session, episode)
	data.InsertMovie(omdbResponse.Search[0], finalFilePath, session, episode)
	writeFinalFileOnDisk(filename, finalFilePath)
}

func respondSeveralMoviesOptions(w http.ResponseWriter, movies omdb.OmdbResponse, filename, session, episode string) {
	w.WriteHeader(200)
	data := getSeveralMoviesOptionsJSON(movies, filename, session, episode)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(data))
}

func getSeveralMoviesOptionsJSON(movies omdb.OmdbResponse, filename, session, episode string) []byte {
	uploadResponse := &SeveralOptionsResponse{Saved: false, Options: movies, Filename: filename, Session: session, Episode: episode}
	data, err := json.Marshal(uploadResponse)
	if err != nil {
		log.Println("Error parsing movie response to JSON:", err)
	}
	return data
}

func getFinalPathForFile(movie *items.Movie, session, episode string) string {
	var finalFilePath string
	if session != "" {
		finalFilePath = finalFilePath + "/" + sessionAb + session + "/"
	}
	finalFilePath = finalFilePath + movie.Title + "(" + movie.Year + ")"
	if episode != "" {
		finalFilePath = finalFilePath + "_" + episodeAb + episode
	}
	finalFilePath = strings.ReplaceAll(finalFilePath, " ", "_")
	return finalFilePath
}

func writeFinalFileOnDisk(input, output string) error {
	in, err := os.Open(tmp + input)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(moviesPath + output)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	err = os.Remove(tmp + input)
	if err != nil {
		return err
	}
	return out.Close()
}
