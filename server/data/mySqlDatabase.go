package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mgarmuno/mediaWebServer/server/items"
)

const (
	dbName             = "./mediaWebServerDatabase.db"
	driver             = "sqlite3"
	create             = "CREATE TABLE IF NOT EXISTS"
	colon              = ":"
	blankSapace        = " "
	interrogation      = "?"
	comma              = ","
	parBeg             = "("
	parEnd             = ")"
	movie              = "movie"
	user               = "user"
	season             = "season"
	actor              = "actor"
	director           = "director"
	movieCast          = "movie_cast"
	movieDirCast       = "movie_dir_cast"
	genre              = "genre"
	movieGenres        = "movie_genres"
	episode            = "episode"
	nullString         = "null"
	integerString      = "integer"
	textString         = "text"
	movieCastFields    = "id integer primary key, id_actor integer, id_movie integer, role text"
	movieDirCastFields = "id integer primary key, id_director integer, id_movie integer"
	movieGenresFields  = "id integer primary key, id_genre integer, id_movie integer"
	moviesDirectory    = "movies_directory"
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
	createTable(database, movie, getObjectFieldsForCreateTabble(&items.Movie{}))
	createTable(database, user, getObjectFieldsForCreateTabble(&items.User{}))
	createTable(database, season, getObjectFieldsForCreateTabble(&items.Season{}))
	createTable(database, actor, getObjectFieldsForCreateTabble(&items.Actor{}))
	createTable(database, director, getObjectFieldsForCreateTabble(&items.Director{}))
	createTable(database, genre, getObjectFieldsForCreateTabble(&items.Genre{}))
	createTable(database, episode, getObjectFieldsForCreateTabble(&items.Episode{}))
	createTable(database, movieCast, strings.Split(movieCastFields, comma))
	createTable(database, movieDirCast, strings.Split(movieDirCastFields, comma))
	createTable(database, movieGenres, strings.Split(movieGenresFields, comma))
	fmt.Println("Database created, closing...")
	database.Close()
}

func createTable(database *sql.DB, table string, fields []string) {
	_, err := database.Exec(create + blankSapace + table + parBeg + strings.Join(fields[:], comma) + parEnd)
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
