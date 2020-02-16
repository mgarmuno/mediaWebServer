package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbName                   = "./mediaWebServerDatabase.db"
	driver                   = "sqlite3"
	create                   = "CREATE TABLE IF NOT EXISTS"
	asterisk                 = "*"
	colon                    = ":"
	blankSapace              = " "
	interrogation            = "?"
	comma                    = ","
	parBeg                   = "("
	parEnd                   = ")"
	argumentForQuery         = "%d"
	movieCreateFields        = "id integer primary key, title text, year text, runtime, text, country text, plot text, rating text, poster text, imdbid text, type text, filepath text"
	userCreateFields         = "id integer primary key, name text, password text"
	actorCreateFields        = "id integer primary key, name text"
	directorCreateFields     = "id integer primary key, name text"
	episodeCreateFields      = "id integer primary key, season_id integer, episode text, filepath text"
	genreCreateFields        = "id integer primary key, name text"
	seasonCreateFields       = "id integer primary key, movie_id integer, season text"
	movieCastCreateFields    = "id integer primary key, id_movie integer, id_actor integer"
	movieDirCastCreateFields = "id integer primary key, id_movie integer, id_director integer"
	movieGenresCreateFields  = "id integer primary key, id_movie integer, id_genre integer"
	movieFields              = "title,year,runtime,country,plot,imdb_rating,poster,imdb_id,type,filepath"
	seasonFields             = "movie_id,season"
	movieTableName           = "movie"
	userTableName            = "user"
	seasonTableName          = "season"
	actorTableName           = "actor"
	directorTableName        = "director"
	movieCastTableName       = "movie_cast"
	movieDirCastTableName    = "movie_dir_cast"
	genreTableName           = "genre"
	movieGenresTableName     = "movie_genres"
	episodeTableName         = "episode"
	nullString               = "null"
	integerString            = "integer"
	textString               = "ext"
	imdbIDFieldName          = "imdb_id"
	movieIDField             = "movie_id"
	seasonNumberField        = "season"
	seasonIDField            = "season_id"
	episodeNumberField       = "episode"
	movieCastFields          = "id integer primary key, id_actor integer, id_movie integer, role text"
	movieDirCastFields       = "id integer primary key, id_director integer, id_movie integer"
	movieGenresFields        = "id integer primary key, id_genre integer, id_movie integer"
	moviesDirectory          = "movies_directory"
)

type DatabaseAPI struct {
}

func CheckDatabase() {
	var databaseCreated bool = true
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		fmt.Println("The database does not exists, creating...")
		databaseCreated = false
	}
	if !databaseCreated {
		database := OpenConnection()
		createDatabaseStructure(&database)
		database.Close()
	}
}

func OpenConnection() sql.DB {
	database, err := sql.Open(driver, dbName)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	return *database
}

func getAppDirectory() string {
	actualDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Cannot access the application directory:", err)
	}
	return actualDir
}

func createDatabaseStructure(database *sql.DB) {
	fmt.Println("Creating database structure...")
	createTable(database, movieTableName, movieFields)
	createTable(database, userTableName, userCreateFields)
	createTable(database, seasonTableName, seasonCreateFields)
	createTable(database, actorTableName, actorCreateFields)
	createTable(database, directorTableName, directorCreateFields)
	createTable(database, genreTableName, genreCreateFields)
	createTable(database, episodeTableName, episodeCreateFields)
	createTable(database, movieCastTableName, movieCastCreateFields)
	createTable(database, movieDirCastTableName, movieDirCastCreateFields)
	createTable(database, movieGenresTableName, movieGenresCreateFields)
	fmt.Println("Database created, closing...")
	database.Close()
}

func createTable(database *sql.DB, table string, fields string) {
	_, err := database.Exec(create + blankSapace + table + parBeg + fields + parEnd)
	if err != nil {
		log.Fatal("Error creating table", table, colon, err)
	}
	fmt.Println("Table", table, "created")
}

func getObjectFieldsForCreateTabble(fields interface{}) []string {
	return getObjectFields(fields, true)
}

func getObjectFields(fields interface{}, isForCreate bool) []string {
	var fieldsSlice []string
	val := reflect.ValueOf(fields).Elem()
	for i := 0; i < val.NumField(); i++ {
		var name string = val.Type().Field(i).Name
		if isForCreate {
			name = addDataTypeToField(val, name, i)
		}
		fieldsSlice = append(fieldsSlice, strings.ToUpper(name))
	}
	return fieldsSlice
}

func getObjectValues(fields interface{}) {

}

func addDataTypeToField(val reflect.Value, name string, i int) string {
	var datatype string = getDatatypeForSQLite(val.Type().Field(i).Type.String())
	if datatype == blankSapace {
		return ""
	}
	var primary string = ""
	if name == "Id" {
		primary = "primary key"
	}
	return name + blankSapace + datatype + blankSapace + primary
}

func getDatatypeForSQLite(goDatatype string) string {
	switch goDatatype {
	case "int":
		return integerString
	case "string":
		return textString
	}
	return blankSapace
}
