package data

import (
	"fmt"
	"strings"

	"github.com/mgarmuno/mediaWebServer/server/items"
)

const (
	insertQuery      = "INSERT INTO"
	movieTable       = "movie"
	movieFieldsTitle = "Title, Year, Runtime, Country, Plot, ImdbRating, Poster, ImdbID, Type, Filepath"
	valuesQuery      = "values"
)

func InsertMovie(movie items.Movie, filepath, session, episode string) {
	insertMovie(movie)
}

func insertMovie(movie items.Movie) {
	database := OpenConnection()
	values := getMovieFieldValues(movie)
	var query string = insertQuery + blankSapace + movieTable + parBeg + movieFieldsTitle + parEnd + valuesQuery + parBeg + values + parEnd
	_, err := database.Exec(query)
	if err != nil {
		fmt.Println("Error inserting into movies:", err)
	}
	fmt.Println("Movie", movie.Title, "inserted")
	database.Close()
}

func getMovieFieldValues(movie items.Movie) string {
	values := []string{
		strings.ReplaceAll(movie.Title, "'", "''"),
		movie.Year,
		movie.Runtime,
		movie.Country,
		strings.ReplaceAll(movie.Plot, "'", "''"),
		movie.ImdbRating,
		movie.Poster,
		movie.ImdbID,
		movie.Type,
		movie.Filepath}
	return "'" + strings.Join(values, "'"+comma+"'") + "'"
}
