package main

import (
	"database/sql"
	"log"

	//I used _ since these libraries are underlined driver
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/jackc/pgx/v4"
)

// We split this code into a function because
// we can connect to different DBs (Mongo, SQL, ...) with different openDB function
// Also, we can create a test db that does not require openDB function

// wrapper function for open Postgres database
func openDB(dsn string) (*sql.DB, error) {
	//use database/sql package Open method
	//which return a db struct and error

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	//use Ping() to check if the db is connected
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// after we make sure that the db is connected, return the database object
	return db, nil
}

// add a method to our application struct to be accessed by our app instance
func (app *application) connectToDB() (*sql.DB, error) {
	//app.DSN is created from the flag that we created and parsed in main.go
	connection, err := openDB(app.DSN)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Postgres!")
	return connection, nil
}
