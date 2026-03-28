package service

import (
	"context"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/repository"
	"github.com/go-telegram/bot"
)

type AdminService struct {
	repo *repository.Repository
	bot  *bot.Bot
}

func NewAdminService(repo *repository.Repository, b *bot.Bot) *AdminService {
	return &AdminService{repo: repo, bot: b}
}

func (s *AdminService) SyncChatAdmins(ctx context.Context, chatID int64) error {
	admins, err := s.bot.GetChatAdministrators(ctx, &bot.GetChatAdministratorsParams{
		ChatID: chatID,
	})
	if err != nil {
		return err
	}

	// clear old
	err = s.repo.DeleteAllAdmins(ctx, chatID)
	if err != nil {
		return err
	}

	for _, a := range admins {
		role := string(a.Type)

		err := s.repo.UpsertAdmin(ctx, int(chatID), a.Member.User.ID, role)

		if err != nil {
			return err
		}
	}

	return nil
}
