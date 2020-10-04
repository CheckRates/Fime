package fime

import "time"

// User of the Fime app
type User struct {
	ID        int64
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time `db:"createdAt"`
}
