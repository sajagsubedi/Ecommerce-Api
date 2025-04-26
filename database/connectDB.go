package database

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func CreateConnection() *sql.DB {

	//connect to db
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URI"))

	if err != nil {
		panic(err)
	}

	//checking connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}
	return db
}
