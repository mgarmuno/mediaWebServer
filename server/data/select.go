package data

import (
	"database/sql"

	"github.com/mgarmuno/mediaWebServer/server/items"
)

const (
	selectTableWhere = "SELECT * FROM %d WHERE"
	equalSign        = "="
)

func GetMovieByImdbID(imdbid string) (items.Movie, error) {
	var database sql.DB = OpenConnection()
	rows, err := database.Query(addAndConditionToQuery(1), movieTableName, imdbIDFieldName, imdbid)
	var movie items.Movie
	if err != nil {
		return movie, err
	}
	for rows.Next() {
		err = rows.Scan(movie)
	}
	rows.Close()
	return movie, nil
}

func GetSeasonByMovieIDAndSeason(movieID int64, seasonNumber string) (items.Season, error) {
	var database sql.DB = OpenConnection()
	rows, err := database.Query(addAndConditionToQuery(2), movieIDField, movieID, seasonNumberField, seasonNumber)
	var season items.Season
	if err != nil {
		return season, err
	}
	for rows.Next() {
		err = rows.Scan(season)
	}
	rows.Close()
	if err != nil {
		return items.Season{}, err
	}
	return season, nil
}

func GetEpisodeBySeasonIDAndEpisode(seasonID int64, episodeNumber string) (items.Episode, error) {
	var database sql.DB = OpenConnection()
	rows, err := database.Query(addAndConditionToQuery(2), seasonIDField, seasonID, episodeNumberField, episodeNumber)
	var episode items.Episode
	if err != nil {
		return episode, err
	}
	for rows.Next() {
		err = rows.Scan(episode)
	}
	rows.Close()
	if err != nil {
		return items.Episode{}, err
	}
	return episode, nil
}

func addAndConditionToQuery(numberOfConditions int) string {
	var query string = selectTableWhere
	for i := 0; i < numberOfConditions; i++ {
		query = query + blankSapace + argumentForQuery + equalSign + argumentForQuery
	}
	return query
}
