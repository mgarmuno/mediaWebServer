package main

import (
	"database/sql"
	"log"
)

const (
	username     = "username"
	password     = "password"
	colom        = ":"
	dat          = "@"
	ip           = "127.0.0.1"
	port         = "3306"
	databaseName = "mediaWebServerDatabase"
	databaseType = "mysql"
)

func openConnection() {
	db, err := sql.Open(databaseType, username+colom+password+dat+"("+ip+colom+port+"/"+databaseName+"parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
}
