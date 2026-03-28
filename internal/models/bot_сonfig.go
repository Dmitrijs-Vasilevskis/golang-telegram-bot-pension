package models

import "time"

type BotConfig struct {
	ChatID             int64      `db:"chat_id"`
	SummaryEnabled     bool       `db:"summary_enabled"`
	DuplicateDMEnabled bool       `db:"duplicate_dm_enabled"`
	UpdatedBy          *int64     `db:"updated_by"`
	CreatedAt          time.Time  `db:"created_at"`
	UpdatedAt          *time.Time `db:"updated_at"`
}
