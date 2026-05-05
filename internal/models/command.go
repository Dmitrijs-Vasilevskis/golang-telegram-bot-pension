package models

type CommandConfig struct {
	ID          int64  `db:"id"`
	ChatID      int64  `db:"chat_id"`
	CommandName string `db:"command_name"`
	Enabled     bool   `db:"enabled"`
}
