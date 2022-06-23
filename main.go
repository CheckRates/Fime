package main

import (
	"log"
	"os"

	"github.com/checkrates/Fime/api"
	"github.com/checkrates/Fime/config"
	"github.com/checkrates/Fime/db/postgres"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	config := config.New()

	// Open & connect to databse
	conn, err := sqlx.Open("postgres", config.Database.ConnString)
	if err != nil {
		log.Fatal("Cannot connect to the database: ", err)
	}

	// Setup Data Access Layer and server
	store, err := postgres.NewStore(conn)
	if err != nil {
		log.Fatal("Cannot create data access store: ", err)
	}

	// Quick fix for Heroku Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	// Start Server
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot create server: ", err)
	}
	err = server.Start(config.Address)
	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}
}
