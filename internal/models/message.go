package models

import "time"

type Message struct {
	ID        int64     `json:"id"`
	ChatID    int64     `json:"chat_id"`
	MessageId int64     `json:"message_id"`
	Username  *string   `json:"username"`
	Text      *string   `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
