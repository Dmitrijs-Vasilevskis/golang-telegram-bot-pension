package models

type Chat struct {
	ID    int64  `db:"id"`
	Title string `db:"title"`
	Type  string `db:"type"`
}
