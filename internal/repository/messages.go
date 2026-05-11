package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/models"
	"github.com/jackc/pgx/v5"
)

const (
	MaxMessagesPerChat = 400
)

func (r *Repository) Save(ctx context.Context, m *models.Message) error {
	query := `
		INSERT INTO messages
		(chat_id, message_id, username, text)
		VALUES (@chat_id, @message_id, @username, @text)
	`

	_, err := r.db.Exec(ctx, query, pgx.NamedArgs{
		"chat_id":    m.ChatID,
		"message_id": m.MessageId,
		"username":   m.Username,
		"text":       m.Text,
	})

	return err
}

func (r *Repository) Update(ctx context.Context, m *models.Message) error {
	query := `
		UPDATE messages
		SET text=@text, updated_at=@updated_at
		WHERE chat_id=@chat_id AND message_id=@message_id
	`
	_, err := r.db.Exec(ctx, query, pgx.NamedArgs{
		"chat_id":    m.ChatID,
		"message_id": m.MessageId,
		"text":       m.Text,
		"updated_at": m.UpdatedAt,
	})

	return err
}

func (r *Repository) GetLastMessages(ctx context.Context, chatID int64, limit int) ([]*models.Message, error) {
	query := `
		SELECT id, chat_id, message_id, username, text, created_at, updated_at
		FROM messages
		WHERE chat_id=$1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, chatID, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	messages, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[models.Message])
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *Repository) TrimMessages(ctx context.Context, chatID int64, limit int) error {
	query := `
		DELETE FROM messages
		WHERE chat_id = $1
		AND id < (
			SELECT id
			FROM messages
			WHERE chat_id = $1
			ORDER BY created_at DESC
			OFFSET $2
			LIMIT 1
		)
	`
	_, err := r.db.Exec(ctx, query, chatID, limit)

	return err
}

func (r *Repository) GetMessagesInRange(ctx context.Context, chatID int64, from, to time.Time) ([]models.Message, error) {
	query := `
		SELECT id, chat_id, message_id, username, text, created_at, updated_at
		FROM messages
		WHERE chat_id = $1
			AND created_at >= $2
			AND created_at <= $3
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, chatID, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages range %w", err)
	}

	defer rows.Close()

	messages, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.Message])
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// unsued
func (r *Repository) CountMessages(ctx context.Context, chatID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM messages WHERE chat_id = $1`

	var count int64
	err := r.db.QueryRow(ctx, query, chatID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count messages: %w", err)
	}

	return count, nil
}

func (r *Repository) DeleteMessagesByChatID(ctx context.Context, chatID int64) (int64, error) {
	query := `DELETE FROM messages WHERE chat_id=$1`

	res, err := r.db.Exec(ctx, query, chatID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete messages for chat %d: %w", chatID, err)
	}

	rowsAffected := res.RowsAffected()

	return rowsAffected, nil
}
