package handlers

import (
	"context"
	"log"
	"time"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/app"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/database"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/helpers"
	messageModel "github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/models"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/repository"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/jackc/pgx/v5"
)

func RecordMessage(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	app *app.App,
) {
	if update.Message == nil {
		return
	}

	message := update.Message
	chat := update.Message.Chat

	if message.Text == "" {
		return
	}

	if message.From == nil || message.From.IsBot {
		return
	}

	if message.ForwardOrigin != nil {
		return
	}

	if helpers.IsUrl(message.Text) {
		return
	}

	currentMessage := &messageModel.Message{
		ChatID:    chat.ID,
		MessageId: int64(message.ID),
		Username:  &message.From.FirstName,
		Text:      &message.Text,
		CreatedAt: time.Now(),
	}

	err := database.WithTransaction(ctx, app.DB, func(tx pgx.Tx) error {
		txRepo := repository.NewRepository(tx)

		if err := txRepo.Save(ctx, currentMessage); err != nil {
			return err
		}

		if err := txRepo.TrimMessages(ctx, chat.ID, repository.MaxMessagesPerChat); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Printf("Failed to record a message from chat: %v, error: %v", chat.ID, err)
	}
}
