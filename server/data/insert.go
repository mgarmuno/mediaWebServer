package data

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/mgarmuno/mediaWebServer/server/items"
)

const (
	insertQuery = "INSERT INTO %d"
	valuesQuery = "values"
)

func InsertMovie(movie items.Movie, filepath, season, episode string) {
	movieFromDB, err := GetMovieByImdbID(movie.ImdbID)
	if err != nil {

	}
	if movieFromDB == (items.Movie{}) {
		var result sql.Result = insertMovie(movie, filepath)
		movieFromDB.ID, _ = result.LastInsertId()
	}
	var seasonID int64 = checkSeasonExists(movieFromDB.ID, season)
	checkEpisodeExists(seasonID, episode)
}

func insertMovie(movie items.Movie, filepath string) sql.Result {
	database := OpenConnection()
	var query string = addInsertFields(10)
	result, err := database.Exec(query, getFieldsForInsert(movieFields, movie))
	if err != nil {
		fmt.Println("Error inserting into movies:", err)
	}
	fmt.Println("Movie", movie.Title, "inserted")
	database.Close()
	return result
}

func getFieldsForInsert(tableFieldNames string, values interface{}) interface{} {
	var fieldsAndValues []string = strings.Split(tableFieldNames, comma)
	structure := reflect.ValueOf(values)
	for i := 0; i < structure.NumField(); i++ {
		if structure.Type().Field(i).Name == "ID" {
			continue
		}
		fieldsAndValues = append(fieldsAndValues, fmt.Sprintf("%v", structure.Field(i).Interface()))
	}
	return fieldsAndValues
}

func checkSeasonExists(movieID int64, season string) int64 {
	seasonFromDB, err := GetSeasonByMovieIDAndSeason(movieID, season)
	if err != nil {
		//TODO
	}
	if seasonFromDB == (items.Season{}) {
		result := insertSeason(int64(movieID), season)
		lastID, _ := result.LastInsertId()
		return lastID
	}
	return seasonFromDB.ID
}

func insertSeason(lastID int64, season string) sql.Result {
	database := OpenConnection()
	var query string = addInsertFields(2)
	insertResult, err := database.Exec(query, seasonTableName, movieIDField, seasonNumberField, lastID, season)
	if err != nil {
		//TODO
	}
	database.Close()
	return insertResult
}

func checkEpisodeExists(seasonID int64, episode string) int64 {
	episodeFromDB, err := GetEpisodeBySeasonIDAndEpisode(seasonID, episode)
	if err != nil {
		//TODO
	}
	if episodeFromDB == (items.Episode{}) {
		result := insertEpisode(int64(seasonID), episode)
		lastID, _ := result.LastInsertId()
		return lastID
	}
	return episodeFromDB.ID
}

func insertEpisode(lastID int64, episode string) sql.Result {
	database := OpenConnection()
	var query string = addInsertFields(2)
	insertResult, err := database.Exec(query, episodeTableName, episodeNumberField, seasonIDField, lastID, episode)
	if err != nil {
		//TODO
	}
	database.Close()
	return insertResult
}

func addInsertFields(numberOfFields int) string {
	var query string = insertQuery + parBeg
	for i := 0; i < numberOfFields; i++ {
		query = query + argumentForQuery + comma
	}
	query = query[:len(query)-1] + parEnd + valuesQuery + parBeg
	for i := 0; i < numberOfFields; i++ {
		query = query + argumentForQuery + comma
	}
	return query[:len(query)-1] + parEnd
}
