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
	dbName                = "./mediaWebServerDatabase.db"
	driver                = "sqlite3"
	create                = "CREATE TABLE IF NOT EXISTS"
	colom                 = ":"
	sp                    = " "
	interrogation         = "?"
	comma                 = ","
	parBeg                = "("
	parEnd                = ")"
	movie                 = "movie"
	serie                 = "serie"
	user                  = "user"
	season                = "season"
	actor                 = "actor"
	director              = "director"
	moviCast              = "movie_cast"
	movieDirCast          = "movie_dir_cast"
	genre                 = "genre"
	movieGenres           = "movie_genres"
	nullString            = "null"
	integerString         = "integer"
	textString            = "text"
	movieCastFields       = "(id_actor integer, id_movie integer, role text, primary key(id_actor, id_movie))"
	movieDirCastFields    = "(id_director integer, id_movie integer, primary key(id_director, id_movie))"
	movieGenresFields     = "(id_genre integer, id_movie integer, primary key(id_genre, id_movie))"
	moviesDirectory       = "movies_directory"
	moviesDirectoryFields = "(id integer primary key, directory text)"
)

type DatabaseAPI struct{}

func OpenConnection() {
	var databaseCreated bool = true
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		fmt.Println("The database does not exists, creating...")
		databaseCreated = false
	}
	database, err := sql.Open(driver, dbName)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	if !databaseCreated {
		createStructure(database)
	}
}

func getAppDirectory() string {
	actualDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Cannot access the application directory:", err)
	}
	return actualDir
}

func createStructure(database *sql.DB) {
	fmt.Println("Creating database structure...")
	createTable(database, movie, getObjectFields(&items.Movie{}))
	createTable(database, serie, getObjectFields(&items.Serie{}))
	createTable(database, user, getObjectFields(&items.User{}))
	createTable(database, season, getObjectFields(&items.Season{}))
	createTable(database, actor, getObjectFields(&items.Actor{}))
	createTable(database, director, getObjectFields(&items.Director{}))
	createTable(database, genre, getObjectFields(&items.Genre{}))
	createTable(database, moviCast, movieCastFields)
	createTable(database, movieDirCast, movieDirCastFields)
	createTable(database, movieGenres, movieGenresFields)
	createTable(database, moviesDirectory, moviesDirectoryFields)
}

func createTable(database *sql.DB, table string, fields string) {
	_, err := database.Exec(create + sp + table + fields)
	if err != nil {
		log.Fatal("Error creating table", table, colom, err)
	}
	fmt.Println("Table", table, "created")
}

func getObjectFields(fields interface{}) string {
	var fieldsSlice []string
	val := reflect.ValueOf(fields).Elem()
	for i := 0; i < val.NumField(); i++ {
		var datatype string = getDatatypeForSqlite(val.Type().Field(i).Type.String())
		if datatype == sp {
			continue
		}
		var name string = val.Type().Field(i).Name
		var primary string = ""
		if name == "id" {
			primary = "primary key"
		}
		fieldsSlice = append(fieldsSlice, name+sp+datatype+sp+primary)
	}
	return parBeg + strings.Join(fieldsSlice, comma) + parEnd
}

func getDatatypeForSqlite(goDatatype string) string {
	switch goDatatype {
	case "int":
		return integerString
	case "string":
		return textString
	case "date":
		return integerString
	}
	return sp
}
