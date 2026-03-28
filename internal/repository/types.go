// internal/repository/types.go
package repository

import "time"

type TelegramAdmin struct {
	UserID    int64
	Username  string
	FirstName string
	LastName  string
	Role      string // "creator" or "administrator"
}

type ChatConfig struct {
	ChatID             int64
	SummaryEnabled     bool
	DuplicateDMEnabled bool
	UpdatedBy          *int64
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
