package main

import (
	"log"

	"github.com/checkrates/Fime/api"
	"github.com/checkrates/Fime/config"
	"github.com/checkrates/Fime/db/postgres"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// TODO: Setup in a env file
const serverAddress = ":3001"

// TODO: Using fime test database before production. Change later
func main() {
	// Open & connect to databse
	conn, err := sqlx.Open("postgres", config.New().Database.ConnString)
	if err != nil {
		log.Fatal("Cannot connect to the database: ", err)
	}

	// Setup Data Access Layer and server
	store, err := postgres.NewStore(conn)
	if err != nil {
		log.Fatal("Cannot create data access store: ", err)
	}

	// Start Server
	server := api.NewServer(store)
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}
}
