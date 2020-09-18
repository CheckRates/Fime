package postgres

import (
	"log"
	"os"
	"testing"

	"github.com/checkrates/Fime/config"
	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	dal, err := NewStore(config.New().Database.ConnString)
	if err != nil {
		log.Fatal("error opening connecting to db: ", err)
	}

	dal.User(1)
	os.Exit(m.Run())
}
