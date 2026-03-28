package repository

import (
	"context"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/models"
)

func (r *Repository) EnsureDefaultCommands(ctx context.Context, chatID int64) error {
	commands := []string{"ask", "summary", "factcheck", "look"}

	for _, cmd := range commands {
		_, err := r.db.Exec(ctx, `
			INSERT INTO command_configs (chat_id, command_name, enabled)
			VALUES ($1, $2, TRUE)
			ON CONFLICT (chat_id, command_name) DO NOTHING
		`, chatID, cmd)

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) ToggleCommand(ctx context.Context, chatID int64, command string) error {
	query := `
		UPDATE command_configs
		SET enabled = NOT enabled
		WHERE chat_id = $1 AND command_name = $2
	`

	_, err := r.db.Exec(ctx, query, chatID, command)
	return err
}

func (r *Repository) SetCommandEnabled(ctx context.Context, chatID int64, command string, enabled bool) error {
	query := `
		UPDATE command_configs
		SET enabled = $1
		WHERE chat_id = $2 AND command_name = $3
	`

	_, err := r.db.Exec(ctx, query, enabled, chatID, command)
	return err
}

func (r *Repository) SetAllCommands(ctx context.Context, chatID int64, enabled bool) error {
	query := `
		UPDATE command_configs
		SET enabled = $1
		WHERE chat_id = $2
	`

	_, err := r.db.Exec(ctx, query, enabled, chatID)
	return err
}

func (r *Repository) GetCommands(ctx context.Context, chatID int64) ([]models.CommandConfig, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, chat_id, command_name, enabled
		FROM command_configs
		WHERE chat_id = $1
	`, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.CommandConfig

	for rows.Next() {
		var c models.CommandConfig
		if err := rows.Scan(&c.ID, &c.ChatID, &c.CommandName, &c.Enabled); err != nil {
			return nil, err
		}
		result = append(result, c)
	}

	return result, nil
}
