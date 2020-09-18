package fime

// Image is a image
type Image struct {
	ID      int64  `db:"id"`
	Name    string `db:"name"`
	URL     string `db:"url"`
	OwnerID int64  `db:"owner"`
}
