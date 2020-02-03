package movie

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type MovieResponse struct {
	Saved    bool
	Filename string
}

func doPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "multipart/form-data")

	file, handler, err := r.FormFile("file")

	if err != nil {
		fmt.Println("Error uploading file:", err)
		return
	}
	defer file.Close()

	osFile, err := os.Create("/tmp/" + handler.Filename)
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		return
	}
	defer osFile.Close()

	hddFile, err := os.OpenFile(osFile.Name(), os.O_RDWR, 0)
	if err != nil {
		fmt.Println("Error opening temporary file:", err)
		return
	}
	defer hddFile.Close()

	_, err = io.Copy(hddFile, file)
	if err != nil {
		fmt.Println("Error writing temporary file:", err)
		return
	}
	fmt.Println("File saved:", handler.Filename)
	w.WriteHeader(200)
	data, err := getMovieSavesResponse(handler.Filename)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(data))
}

func getMovieSavesResponse(filename string) ([]byte, error) {
	movieResponse := &MovieResponse{Saved: true, Filename: filename}
	data, err := json.Marshal(movieResponse)
	return data, err
}
