package repository

import (
	"context"
	"fmt"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/models"
	"github.com/jackc/pgx/v5"
)

// RegisterChatWithDefaults registers a new chat with default settings
func (r *Repository) RegisterChatWithDefaults(ctx context.Context, chatID int64, title, chatType string, addedByUserID int64, addedByUsername, addedByFirstName, addedByLastName string) error {
	tx, err := r.beginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start tx for chat registration: %w", err)
	}
	defer tx.Rollback(ctx)

	txRepo := NewRepository(tx)

	// Create a user who added the bot
	_, err = txRepo.db.Exec(ctx, `
		INSERT INTO users (id, username, first_name, last_name, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (id) DO UPDATE SET
			username = EXCLUDED.username,
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name
	`, addedByUserID, addedByUsername, addedByFirstName, addedByLastName)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Create a chat
	_, err = txRepo.db.Exec(ctx, `
		INSERT INTO chats (id, title, type, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			type = EXCLUDED.type
	`, chatID, title, chatType)
	if err != nil {
		return fmt.Errorf("failed to create chat: %w", err)
	}

	// Add a user as an administrator
	_, err = txRepo.db.Exec(ctx, `
		INSERT INTO chat_admins (chat_id, user_id, role)
		VALUES ($1, $2, 'administrator')
		ON CONFLICT (chat_id, user_id) DO NOTHING
	`, chatID, addedByUserID)
	if err != nil {
		return fmt.Errorf("failed to add admin: %w", err)
	}

	// Create a bot configuration (summary and duplicate are disabled by default)
	_, err = txRepo.db.Exec(ctx, `
		INSERT INTO bot_configs (chat_id, summary_enabled, duplicate_dm_enabled, updated_by, created_at, updated_at)
		VALUES ($1, false, false, $2, NOW(), NOW())
		ON CONFLICT (chat_id) DO UPDATE SET
			summary_enabled = EXCLUDED.summary_enabled,
			duplicate_dm_enabled = EXCLUDED.duplicate_dm_enabled,
			updated_by = EXCLUDED.updated_by,
			updated_at = NOW()
	`, chatID, addedByUserID)
	if err != nil {
		return fmt.Errorf("failed to create bot config: %w", err)
	}

	// Create records for all commands (enabled by default)
	_, err = txRepo.db.Exec(ctx, `
		INSERT INTO command_configs (chat_id, command_name, enabled)
		VALUES
			($1, 'ask', true),
			($1, 'summary', false),
			($1, 'factcheck', false),
			($1, 'look', true),
			($1, 'status', true)
		ON CONFLICT (chat_id, command_name) DO NOTHING
	`, chatID)
	if err != nil {
		return fmt.Errorf("failed to create command configs: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit chat registration tx: %w", err)
	}

	return nil
}

// SyncChatAdmins syncs chat administrators
func (r *Repository) SyncChatAdmins(ctx context.Context, chatID int64, admins []TelegramAdmin) error {
	tx, err := r.beginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start tx for admin sync: %w", err)
	}
	defer tx.Rollback(ctx)

	// Removing old administrators
	_, err = tx.Exec(ctx, `DELETE FROM chat_admins WHERE chat_id = $1`, chatID)
	if err != nil {
		return fmt.Errorf("failed to delete old admins: %w", err)
	}

	// Adding new administrators
	for _, admin := range admins {
		// At the first create a user
		_, err = tx.Exec(ctx, `
			INSERT INTO users (id, username, first_name, last_name, created_at)
			VALUES ($1, $2, $3, $4, NOW())
			ON CONFLICT (id) DO UPDATE SET
				username = EXCLUDED.username,
				first_name = EXCLUDED.first_name,
				last_name = EXCLUDED.last_name
		`, admin.UserID, admin.Username, admin.FirstName, admin.LastName)
		if err != nil {
			return fmt.Errorf("failed to create user %d: %w", admin.UserID, err)
		}

		// sets as admin/creator
		_, err = tx.Exec(ctx, `
			INSERT INTO chat_admins (chat_id, user_id, role)
			VALUES ($1, $2, $3)
		`, chatID, admin.UserID, admin.Role)
		if err != nil {
			return fmt.Errorf("failed to add admin %d: %w", admin.UserID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit admin sync tx: %w", err)
	}

	return nil
}

func (r *Repository) DeleteChatById(ctx context.Context, chatID int64) (bool, error) {
	cmdTag, err := r.db.Exec(ctx, `DELETE FROM chats WHERE id = $1`, chatID)

	if err != nil {
		return false, fmt.Errorf("failed to delete chat: %w", err)
	}

	return cmdTag.RowsAffected() == 1, nil
}

func (r *Repository) GetChatConfig(ctx context.Context, chatID int64) (*models.BotConfig, error) {
	var config models.BotConfig
	err := r.db.QueryRow(ctx, `
		SELECT chat_id, summary_enabled, duplicate_dm_enabled, updated_by, created_at, updated_at
		FROM bot_configs
		WHERE chat_id = $1
	`, chatID).Scan(
		&config.ChatID,
		&config.SummaryEnabled,
		&config.DuplicateDMEnabled,
		&config.UpdatedBy,
		&config.CreatedAt,
		&config.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat config: %w", err)
	}
	return &config, nil
}

func (r *Repository) GetCommandConfigs(ctx context.Context, chatID int64) (map[string]bool, error) {
	rows, err := r.db.Query(ctx, `
		SELECT command_name, enabled
		FROM command_configs
		WHERE chat_id = $1
	`, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get command configs: %w", err)
	}
	defer rows.Close()

	configs := make(map[string]bool)
	for rows.Next() {
		var cmd string
		var enabled bool
		if err := rows.Scan(&cmd, &enabled); err != nil {
			return nil, err
		}
		configs[cmd] = enabled
	}

	return configs, nil
}

func (r *Repository) GetFeatureConfig(ctx context.Context, chatID int64, feature string) (bool, error) {
	var col string

	switch feature {
	case "summary":
		col = "summary_enabled"
	case "duplicate_dm":
		col = "duplicate_dm_enabled"
	default:
		return false, fmt.Errorf("unknown feature: %s", feature)
	}

	query := fmt.Sprintf(`
		SELECT %s
		FROM bot_configs
		WHERE chat_id = $1
	`, col)

	var enabled bool
	err := r.db.QueryRow(ctx, query, chatID).Scan(&enabled)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, ErrNotFound
		}
		return false, fmt.Errorf("failed to get feature config: %w", err)
	}

	return enabled, nil
}

func (r *Repository) SetFeature(ctx context.Context, chatID int64, userID int64, feature string, enabled bool) error {
	var col string

	switch feature {
	case "summary":
		col = "summary_enabled"
	case "duplicate_dm":
		col = "duplicate_dm_enabled"
	default:
		return fmt.Errorf("unknown feature: %s", feature)
	}

	query := fmt.Sprintf(`
		UPDATE bot_configs
		SET %s = $1,
			updated_by = $2,
			updated_at = NOW()
			WHERE chat_id = $3
	`, col)

	tag, err := r.db.Exec(ctx, query, enabled, userID, chatID)
	if err != nil {
		return fmt.Errorf("Failed to update feature %s: %w", feature, err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *Repository) beginTx(ctx context.Context) (pgx.Tx, error) {
	txStarter, ok := r.db.(interface {
		Begin(context.Context) (pgx.Tx, error)
	})
	if !ok {
		return nil, fmt.Errorf("database connection does not support transactions")
	}

	return txStarter.Begin(ctx)
}
