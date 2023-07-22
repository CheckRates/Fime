package main

import (
	"log"

	"github.com/checkrates/Fime/config"
	"github.com/checkrates/Fime/rest"
	_ "github.com/lib/pq"
)

func main() {
	config, err := config.Load(".")
	if err != nil {
		log.Fatal("Cannot read config: ", err)
	}

	server, err := rest.NewServer(config)
	if err != nil {
		log.Fatal("Cannot create server: ", err)
	}

	err = server.Start(config.Address)
	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}
}
