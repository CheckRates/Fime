package fime

// Tag image tag
type Tag struct {
	ID   int64  `db:"id"`
	Name string `db:"tag"`
}
