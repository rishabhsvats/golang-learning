package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

func main() {

	//connect to the database
	db := initDB()
	db.Ping()
	// create sessions

	//create channels

	//create waitgroup

	//set up the application config

	//setup mail

	//listen for web connections
}

func initDB() *sql.DB {

	conn := connectToDB()
	if conn != nil {
		log.Panic("can't connect to the database")
	}
	return conn
}

func connectToDB() *sql.DB {
	counts := 0
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("postgres not yet ready")
		} else {
			log.Println("connect to database!")
			return connection
		}

		if counts > 10 {
			return nil
		}
		log.Println("backing off for 1 second")
		time.Sleep(1 * time.Second)
		counts++
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
