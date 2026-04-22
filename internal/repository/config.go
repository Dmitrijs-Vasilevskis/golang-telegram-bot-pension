package repository

import (
	"context"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/models"
	"github.com/jackc/pgx/v5"
)

// unsued
func (r *Repository) EnsureBotConfig(ctx context.Context, chatID int64) error {
	query := `
		INSERT INTO bot_configs (chat_id, summary_enabled, duplicate_dm_enabled)
		VALUES ($1, FALSE, FALSE)
		ON CONFLICT (chat_id) DO NOTHING
	`

	_, err := r.db.Exec(ctx, query, chatID)
	return err
}

// unsued
func (r *Repository) GetBotConfig(ctx context.Context, chatID int64) (*models.BotConfig, error) {
	query := `
		SELECT chat_id, summary_enabled, duplicate_dm_enabled, updated_by
		FROM bot_configs
		WHERE chat_id = $1
	`

	row := r.db.QueryRow(ctx, query, chatID)

	var cfg models.BotConfig
	err := row.Scan(
		&cfg.ChatID,
		&cfg.SummaryEnabled,
		&cfg.DuplicateDMEnabled,
		&cfg.UpdatedBy,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &cfg, nil
}

// unsued
func (r *Repository) SetSummary(ctx context.Context, chatID int64, enabled bool, userID int64) error {
	query := `
		UPDATE bot_configs
		SET summary_enabled = $1, updated_by = $2
		WHERE chat_id = $3
	`

	tag, err := r.db.Exec(ctx, query, enabled, userID, chatID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// unsued
func (r *Repository) SetDuplicateDM(ctx context.Context, chatID int64, enabled bool, userID int64) error {
	query := `
		UPDATE bot_configs
		SET duplicate_dm_enabled = $1, updated_by = $2
		WHERE chat_id = $3
	`

	tag, err := r.db.Exec(ctx, query, enabled, userID, chatID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
