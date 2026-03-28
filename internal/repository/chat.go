// internal/repository/chat.go
package repository

import (
	"context"
	"fmt"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/models"
)

// EnsureUserExists создаёт или обновляет пользователя
func (r *Repository) EnsureUserExists(ctx context.Context, userID int64, username, firstName, lastName string) error {
	query := `
		INSERT INTO users (id, username, first_name, last_name, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (id) DO UPDATE SET
			username = EXCLUDED.username,
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name
	`

	_, err := r.db.Exec(ctx, query, userID, username, firstName, lastName)
	if err != nil {
		return fmt.Errorf("failed to ensure user exists: %w", err)
	}
	return nil
}

// EnsureChatExists создаёт запись о чате, если её нет
func (r *Repository) EnsureChatExists(ctx context.Context, chatID int64, title, chatType string) error {
	query := `
		INSERT INTO chats (id, title, type, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			type = EXCLUDED.type
	`

	_, err := r.db.Exec(ctx, query, chatID, title, chatType)
	if err != nil {
		return fmt.Errorf("failed to ensure chat exists: %w", err)
	}
	return nil
}

// RegisterChatWithDefaults регистрирует новый чат с настройками по умолчанию
func (r *Repository) RegisterChatWithDefaults(ctx context.Context, chatID int64, title, chatType string, addedByUserID int64, addedByUsername, addedByFirstName, addedByLastName string) error {

	// 1. Создаём пользователя, который добавил бота
	_, err := r.db.Exec(ctx, `
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

	// 2. Создаём чат
	_, err = r.db.Exec(ctx, `
		INSERT INTO chats (id, title, type, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			type = EXCLUDED.type
	`, chatID, title, chatType)
	if err != nil {
		return fmt.Errorf("failed to create chat: %w", err)
	}

	// 3. Добавляем пользователя как администратора
	_, err = r.db.Exec(ctx, `
		INSERT INTO chat_admins (chat_id, user_id, role)
		VALUES ($1, $2, 'administrator')
		ON CONFLICT (chat_id, user_id) DO NOTHING
	`, chatID, addedByUserID)
	if err != nil {
		return fmt.Errorf("failed to add admin: %w", err)
	}

	// 4. Создаём конфигурацию бота (summary и duplicate отключены по умолчанию)
	_, err = r.db.Exec(ctx, `
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

	// 5. Создаём записи для всех команд (включены по умолчанию)
	commands := []string{"ask", "summary", "factcheck", "look", "status"}
	for _, cmd := range commands {
		_, err = r.db.Exec(ctx, `
			INSERT INTO command_configs (chat_id, command_name, enabled)
			VALUES ($1, $2, true)
			ON CONFLICT (chat_id, command_name) DO NOTHING
		`, chatID, cmd)
		if err != nil {
			return fmt.Errorf("failed to create command config for %s: %w", cmd, err)
		}
	}

	return nil
}

// SyncChatAdmins синхронизирует администраторов чата
func (r *Repository) SyncChatAdmins(ctx context.Context, chatID int64, admins []TelegramAdmin) error {

	// 1. Удаляем старых администраторов
	_, err := r.db.Exec(ctx, `DELETE FROM chat_admins WHERE chat_id = $1`, chatID)
	if err != nil {
		return fmt.Errorf("failed to delete old admins: %w", err)
	}

	// 2. Добавляем новых администраторов
	for _, admin := range admins {
		// Сначала создаём пользователя
		_, err = r.db.Exec(ctx, `
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

		// Добавляем как администратора
		_, err = r.db.Exec(ctx, `
			INSERT INTO chat_admins (chat_id, user_id, role)
			VALUES ($1, $2, $3)
		`, chatID, admin.UserID, admin.Role)
		if err != nil {
			return fmt.Errorf("failed to add admin %d: %w", admin.UserID, err)
		}
	}

	return nil
}

// GetChatConfig возвращает конфигурацию чата
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

// GetCommandConfigs возвращает настройки команд для чата
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
