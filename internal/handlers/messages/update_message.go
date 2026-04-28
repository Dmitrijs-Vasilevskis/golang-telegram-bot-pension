package handlers

import (
	"context"
	"log"
	"time"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/database"
	moderHandler "github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/handlers/moderation"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/helpers"
	messageModel "github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/models"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/repository"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UpdateMessage(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	db *pgxpool.Pool,
) {
	message := update.EditedMessage

	if helpers.IsToxic(message.Text) {
		moderHandler.Clown(ctx, b, message)
	}

	now := time.Now()
	updatedMessage := &messageModel.Message{
		ChatID:    message.Chat.ID,
		MessageId: int64(message.ID),
		Text:      &message.Text,
		UpdatedAt: &now,
	}

	err := database.WithTransaction(ctx, db, func(tx pgx.Tx) error {
		txRepo := repository.NewRepository(tx)

		if err := txRepo.Update(ctx, updatedMessage); err != nil {
			return err
		}

		if err := txRepo.TrimMessages(ctx, int64(message.ID), repository.MaxMessagesPerChat); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Printf("Failed to update edited message from chat: %v, error: %v", message.Chat.ID, err)
	}
}
