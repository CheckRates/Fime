package postgres

import (
	"log"
	"os"
	"testing"

	"github.com/checkrates/Fime/config"
	_ "github.com/lib/pq"
)

var dal *Store

func init() {
	var err error
	dal, err = NewStore(config.New().Database.ConnString)
	if err != nil {
		log.Fatal("error opening connecting to db: ", err)
	}
}

func TestMain(m *testing.M) {
	dal.User(1)
	os.Exit(m.Run())
}
