// internal/service/chat.go
package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/logger"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/repository"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

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

type ChatService struct {
	repo *repository.Repository
	bot  *bot.Bot
}

func NewChatService(repo *repository.Repository, b *bot.Bot) *ChatService {
	return &ChatService{repo: repo, bot: b}
}

// RegisterChat регистрирует чат в системе при добавлении бота
func (s *ChatService) RegisterChat(ctx context.Context, chatID int64, title, chatType string, addedByUser *models.User) error {
	// 1. Регистрируем чат с настройками по умолчанию
	err := s.repo.RegisterChatWithDefaults(
		ctx,
		chatID,
		title,
		chatType,
		addedByUser.ID,
		addedByUser.Username,
		addedByUser.FirstName,
		addedByUser.LastName,
	)
	if err != nil {
		return fmt.Errorf("failed to register chat: %w", err)
	}

	// 2. Синхронизируем всех администраторов чата
	if err := s.SyncChatAdmins(ctx, chatID); err != nil {
		// Логируем ошибку, но не прерываем регистрацию
		log.Printf("Warning: failed to sync admins for chat %d: %v", chatID, err)
	}

	log.Printf("Chat %d (%s) successfully registered by user %d", chatID, title, addedByUser.ID)
	return nil
}

// SyncChatAdmins синхронизирует всех администраторов чата из Telegram
func (s *ChatService) SyncChatAdmins(ctx context.Context, chatID int64) error {
	// Получаем администраторов из Telegram
	admins, err := s.bot.GetChatAdministrators(ctx, &bot.GetChatAdministratorsParams{
		ChatID: chatID,
	})
	if err != nil {
		return fmt.Errorf("failed to get admins from Telegram: %w", err)
	}

	logger.DebugJson(admins)

	// Конвертируем в наш формат
	var repoAdmins []repository.TelegramAdmin
	for _, admin := range admins {
		role := string(admin.Type)

		switch role {
		case "administrator":
			repoAdmins = append(repoAdmins, repository.TelegramAdmin{
				UserID:    admin.Administrator.User.ID,
				Username:  admin.Administrator.User.Username,
				FirstName: admin.Administrator.User.FirstName,
				LastName:  admin.Administrator.User.LastName,
				Role:      role,
			})
		case "creator":
			repoAdmins = append(repoAdmins, repository.TelegramAdmin{
				UserID:    admin.Owner.User.ID,
				Username:  admin.Owner.User.Username,
				FirstName: admin.Owner.User.FirstName,
				LastName:  admin.Owner.User.LastName,
				Role:      role,
			})
		}
	}

	// Сохраняем в базу
	if err := s.repo.SyncChatAdmins(ctx, chatID, repoAdmins); err != nil {
		return fmt.Errorf("failed to sync admins to DB: %w", err)
	}

	return nil
}

// IsCommandEnabled проверяет, включена ли команда в чате
func (s *ChatService) IsCommandEnabled(ctx context.Context, chatID int64, command string) (bool, error) {
	configs, err := s.repo.GetCommandConfigs(ctx, chatID)
	if err != nil {
		// Если нет конфигурации, команда включена по умолчанию
		return true, nil
	}

	enabled, exists := configs[command]
	if !exists {
		return true, nil
	}

	return enabled, nil
}

// IsSummaryEnabled проверяет, включена ли функция summary
func (s *ChatService) IsSummaryEnabled(ctx context.Context, chatID int64) (bool, error) {
	config, err := s.repo.GetChatConfig(ctx, chatID)
	if err != nil {
		// Если нет конфигурации, функция выключена по умолчанию
		return false, nil
	}
	return config.SummaryEnabled, nil
}

// IsDuplicateDMEnabled проверяет, включено ли дублирование DM
func (s *ChatService) IsDuplicateDMEnabled(ctx context.Context, chatID int64) (bool, error) {
	config, err := s.repo.GetChatConfig(ctx, chatID)
	if err != nil {
		return false, nil
	}
	return config.DuplicateDMEnabled, nil
}

// internal/service/chat.go
func (s *ChatService) GetUserChats(ctx context.Context, userID int64) ([]models.Chat, error) {
	// Получаем ID чатов, где пользователь является администратором
	adminChats, err := s.repo.GetUserAdminChats(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Получаем информацию о чатах из Telegram
	var chats []models.Chat
	for _, adminChat := range adminChats {
		chat, err := s.bot.GetChat(ctx, &bot.GetChatParams{
			ChatID: adminChat.ChatID,
		})
		if err != nil {
			// Если не можем получить чат, пропускаем
			continue
		}
		formattedChat := &models.Chat{
			ID:               chat.ID,
			Type:             chat.Type,
			Title:            chat.Title,
			Username:         chat.Username,
			FirstName:        chat.FirstName,
			LastName:         chat.LastName,
			IsForum:          chat.IsForum,
			IsDirectMessages: chat.IsDirectMessages,
		}

		chats = append(chats, *formattedChat)
	}

	return chats, nil
}
