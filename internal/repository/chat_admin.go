package repository

import (
	"context"
	"fmt"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/models"
)

func (r *Repository) IsAdmin(ctx context.Context, chatID, userID int64) (bool, error) {
	var exists bool

	query := `
		SELECT EXISTS(
			SELECT 1 FROM chat_admins
			WHERE chat_id = $1 AND user_id = $2
		)
	`

	err := r.db.QueryRow(ctx, query, chatID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *Repository) UpsertAdmin(ctx context.Context, chatID int, userID int64, role string) error {
	query := `
		INSERT INTO chat_admins (chat_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (chat_id, user_id) DO UPDATE SET
			role = EXCLUDED.role
	`

	_, err := r.db.Exec(ctx, query, chatID, userID, role)
	if err != nil {
		return fmt.Errorf("failed to upsert admin %d for chat %d: %w", userID, chatID, err)
	}

	return nil
}

func (r *Repository) DeleteAllAdmins(ctx context.Context, chatID int64) error {
	query := `DELETE FROM chat_admins WHERE chat_id = $1`

	_, err := r.db.Exec(ctx, query, chatID)
	if err != nil {
		return fmt.Errorf("failed to delete admins for chat %d: %w", chatID, err)
	}

	return nil
}

func (r *Repository) GetUserAdminChats(ctx context.Context, userID int64) ([]models.ChatAdmin, error) {
	rows, err := r.db.Query(ctx, `
	SELECT chat_id, user_id, role
	FROM chat_admins
	WHERE user_id = $1
	ORDER BY chat_id
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user admin chats: %w", err)
	}

	defer rows.Close()

	var admins []models.ChatAdmin

	for rows.Next() {
		var admin models.ChatAdmin
		if err := rows.Scan(&admin.ChatID, &admin.UserID, &admin.Role); err != nil {
			return nil, err
		}

		admins = append(admins, admin)
	}

	return admins, nil
}
