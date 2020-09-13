package fime

import "time"

// User of the Fime app
type User struct {
	ID        int64
	Name      string
	CreatedAt time.Time `db:"createdAt"`
}
