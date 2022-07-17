package postgres

import (
	"log"
	"os"
	"testing"

	"github.com/checkrates/Fime/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var post *PostSQL
var image *ImageSQL
var user *UserSQL
var tag *TagSQL
var imageTag *ImageTagSQL

func TestMain(m *testing.M) {
	// Load config and test database connection
	config, err := config.Load("../../")
	if err != nil {
		log.Fatal("TEST: Cannot load configuration file")
	}
	conn, err := sqlx.Open("postgres", config.ConnString)
	if err != nil {
		log.Fatal("TEST: Cannot connect to the database: ", err)
	}

	// Initiate all postgres repositories
	post = NewPostRepository(conn)
	image = NewImageRepository(conn)
	user = NewUserRepository(conn)
	tag = NewTagRepository(conn)
	imageTag = NewImageTagRepository(conn)

	os.Exit(m.Run())
}
