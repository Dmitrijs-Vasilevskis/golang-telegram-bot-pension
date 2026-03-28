package service

import (
	"context"
	"fmt"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/repository"
)

type ConfigService struct {
	repo *repository.Repository
}

func NewConfigService(repo *repository.Repository) *ConfigService {
	return &ConfigService{repo: repo}
}

func (s *ConfigService) ensureAdmin(ctx context.Context, chatID, userID int64) error {
	ok, err := s.repo.IsAdmin(ctx, chatID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("not authorized")
	}
	return nil
}

func (s *ConfigService) SetSummary(ctx context.Context, chatID, userID int64, enabled bool) error {
	if err := s.ensureAdmin(ctx, chatID, userID); err != nil {
		return err
	}

	if err := s.repo.EnsureBotConfig(ctx, chatID); err != nil {
		return err
	}

	return s.repo.SetSummary(ctx, chatID, enabled, userID)
}

func (s *ConfigService) SetDuplicateDM(ctx context.Context, chatID, userID int64, enabled bool) error {
	if err := s.ensureAdmin(ctx, chatID, userID); err != nil {
		return err
	}

	if err := s.repo.EnsureBotConfig(ctx, chatID); err != nil {
		return err
	}

	return s.repo.SetDuplicateDM(ctx, chatID, enabled, userID)
}

func (s *ConfigService) ToggleCommand(ctx context.Context, chatID, userID int64, cmd string) error {
	if err := s.ensureAdmin(ctx, chatID, userID); err != nil {
		return err
	}

	if err := s.repo.EnsureDefaultCommands(ctx, chatID); err != nil {
		return err
	}

	return s.repo.ToggleCommand(ctx, chatID, cmd)
}

func (s *ConfigService) SetCommandEnabled(ctx context.Context, chatID, userID int64, cmd string, enabled bool) error {
	if err := s.ensureAdmin(ctx, chatID, userID); err != nil {
		return err
	}

	if err := s.repo.EnsureDefaultCommands(ctx, chatID); err != nil {
		return err
	}

	return s.repo.SetCommandEnabled(ctx, chatID, cmd, enabled)
}

func (s *ConfigService) SetAllCommands(ctx context.Context, chatID, userID int64, enabled bool) error {
	if err := s.ensureAdmin(ctx, chatID, userID); err != nil {
		return err
	}

	if err := s.repo.EnsureDefaultCommands(ctx, chatID); err != nil {
		return err
	}

	return s.repo.SetAllCommands(ctx, chatID, enabled)
}
