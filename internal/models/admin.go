package models

type ChatAdmin struct {
	ChatID int64  `db:"chat_id"`
	UserID int64  `db:"user_id"`
	Role   string `db:"role"` // creator | administrator
}
