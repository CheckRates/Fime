package main

import (
	"fmt"

	"github.com/checkrates/Fime/config"
	"github.com/checkrates/Fime/db/postgres"
	"github.com/checkrates/Fime/fime"
	_ "github.com/lib/pq"
)

func main() {
	config := config.New()
	db, _ := postgres.NewStore(config.Database.ConnString)

	user := fime.User{
		Name: "TestCool",
	}

	db.CreateUser(&user)
	fmt.Println(user.ID)
	fmt.Println(user.CreatedAt)
	fmt.Println(user.Name)
}
