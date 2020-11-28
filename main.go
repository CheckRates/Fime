package main

import (
	"log"

	"github.com/checkrates/Fime/api"
	"github.com/checkrates/Fime/config"
	"github.com/checkrates/Fime/db/postgres"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const serverAddress = ":8000"

// TODO: Using fime test database before production. Change later
func main() {
	conn, err := sqlx.Open("postgres", config.New().Database.ConnString)
	if err != nil {
		log.Fatal("Cannot connect to the database: ", err)
	}

	// Setup Data Access Layer and server
	store, err := postgres.NewStore(conn)
	if err != nil {
		log.Fatal("Cannot create data access store: ", err)
	}

	server := api.NewServer(store)
	// Start Server
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}

}
