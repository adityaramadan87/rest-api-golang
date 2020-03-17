package database

import (
	"fmt"
	_ "log"
	"database/sql"
	_ "github.com/lib/pq"
)
const (
	host = "localhost"
	port = 5432
	user = "postgres"
	password = "terlalupendek"
	dbname = "belajargo"
)


func Connect() *sql.DB  {
	// var username string = "postgres"
	// var host string = "localhost"
	// var password string = "terlalupendek"
	// var database string = "belajargo"
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)


	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}
	// defer db.Close()

	err = db.Ping()
	if err != nil{
		panic(err)
	}

	fmt.Println("Connected succesfully!")
	return db
}

