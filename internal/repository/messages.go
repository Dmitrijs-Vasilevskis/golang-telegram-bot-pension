package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepository struct {
	db *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) *MessageRepository {
	return &MessageRepository{db}
}

func (r *MessageRepository) Save(ctx context.Context, m *models.Message) error {
	query := `
	INSERT INTO messages
	(chat_id, message_id, username, text)
	VALUES ($1,$2,$3,$4)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		m.ChatID,
		m.MessageId,
		m.Username,
		m.Text,
	)

	return err
}

func (r *MessageRepository) GetLastMessages(ctx context.Context, chatID int64, limit int) ([]models.Message, error) {
	query := `
	SELECT id, message_id, username, text, created_at
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

	var messages []models.Message

	for rows.Next() {
		var m models.Message

		err := rows.Scan(
			&m.ID,
			&m.MessageId,
			&m.Username,
			&m.Text,
			&m.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		messages = append(messages, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *MessageRepository) TrimMessages(ctx context.Context, chatID int64, limit int) error {
	query := `
		DELETE FROM messages
		WHERE id IN(
			SELECT id
			FROM messages
			WHERE chat_id=$1
			ORDER BY created_at DESC
			OFFSET $2
		)
	`
	_, err := r.db.Exec(ctx, query, chatID, limit)

	return err
}

func (r *MessageRepository) GetMessagesInRange(ctx context.Context, chatID int64, from, to time.Time) ([]models.Message, error) {
	query := `
		SELECT id, chat_id, username, text, created_at
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

	var messages []models.Message
	for rows.Next() {
		var m models.Message
		err := rows.Scan(
			&m.ID,
			&m.ChatID,
			&m.Username,
			&m.Text,
			&m.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message %w", err)
		}

		messages = append(messages, m)
	}

	return messages, nil
}

func (r *MessageRepository) CountMessages(ctx context.Context, chatID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM messages WHERE chat_d = $1`

	var count int64
	err := r.db.QueryRow(ctx, query, chatID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count messages: %w", err)
	}

	return count, nil
}
