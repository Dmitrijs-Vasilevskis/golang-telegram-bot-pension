package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/logger"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/repository"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type ChatService struct {
	repo *repository.Repository
	bot  *bot.Bot
}

func NewChatService(repo *repository.Repository, b *bot.Bot) *ChatService {
	return &ChatService{repo: repo, bot: b}
}

func (s *ChatService) RegisterChat(ctx context.Context, chatID int64, title, chatType string, addedByUser *models.User) error {
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

	if err := s.SyncChatAdmins(ctx, chatID); err != nil {
		log.Printf("Warning: failed to sync admins for chat %d: %v", chatID, err)
	}

	log.Printf("Chat %d (%s) successfully registered by user %d", chatID, title, addedByUser.ID)
	return nil
}

func (s *ChatService) SyncChatAdmins(ctx context.Context, chatID int64) error {
	admins, err := s.bot.GetChatAdministrators(ctx, &bot.GetChatAdministratorsParams{
		ChatID: chatID,
	})
	if err != nil {
		return fmt.Errorf("failed to get admins from Telegram: %w", err)
	}

	logger.DebugJson(admins)

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

	if err := s.repo.SyncChatAdmins(ctx, chatID, repoAdmins); err != nil {
		return fmt.Errorf("failed to sync admins to DB: %w", err)
	}

	return nil
}

func (s *ChatService) IsCommandEnabled(ctx context.Context, chatID int64, command string) (bool, error) {
	configs, err := s.repo.GetCommandConfigs(ctx, chatID)
	if err != nil {
		return false, err
	}

	enabled, exists := configs[command]
	if !exists {
		return true, nil
	}

	return enabled, nil
}

func (s *ChatService) IsFeatureEnabled(ctx context.Context, chatID int64, feature string) (bool, error) {
	cfg, err := s.repo.GetFeatureConfig(ctx, chatID, feature)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return false, nil
		}
		return false, err
	}

	return cfg, nil
}

// unsued
func (s *ChatService) IsSummaryEnabled(ctx context.Context, chatID int64) (bool, error) {
	config, err := s.repo.GetChatConfig(ctx, chatID)
	if err != nil {
		// Если нет конфигурации, функция выключена по умолчанию
		return false, nil
	}
	return config.SummaryEnabled, nil
}

// unsued
func (s *ChatService) IsDuplicateDMEnabled(ctx context.Context, chatID int64) (bool, error) {
	config, err := s.repo.GetChatConfig(ctx, chatID)
	if err != nil {
		return false, nil
	}
	return config.DuplicateDMEnabled, nil
}

func (s *ChatService) GetUserChats(ctx context.Context, userID int64) ([]models.Chat, error) {
	adminChats, err := s.repo.GetUserAdminChats(ctx, userID)
	if err != nil {
		return nil, err
	}

	var chats []models.Chat
	for _, adminChat := range adminChats {
		chat, err := s.bot.GetChat(ctx, &bot.GetChatParams{
			ChatID: adminChat.ChatID,
		})
		if err != nil {
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

func (s *ChatService) SetFeature(ctx context.Context, chatID int64, userID int64, feature string, enabled bool) error {
	return s.repo.SetFeature(ctx, chatID, userID, feature, enabled)
}
