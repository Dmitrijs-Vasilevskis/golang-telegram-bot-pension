package handlers

import (
	"context"
	"log"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/database"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/repository"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func HandleLeaveChat(ctx context.Context, b *bot.Bot, update *models.Update, db *pgxpool.Pool) {
	if update == nil || update.MyChatMember == nil {
		log.Println("Error: Invalid update in HandleLeaveChat")
		return
	}

	botID := b.ID()
	chatMember := update.MyChatMember

	var isBotLeaving bool
	var chatID int64

	if chatMember.NewChatMember.Left != nil &&
		chatMember.NewChatMember.Left.User != nil {

		leftUser := chatMember.NewChatMember.Left.User

		if leftUser.ID == botID {
			isBotLeaving = true
			chatID = chatMember.Chat.ID
		}
	}

	if !isBotLeaving {
		return
	}

	log.Printf("Bot leaving chat %d, cleaning up messages...", chatID)

	err := database.WithTransaction(ctx, db, func(tx pgx.Tx) error {
		repo := repository.NewRepository(tx)

		deletedCount, err := repo.DeleteMessagesByChatID(ctx, chatID)
		if err != nil {
			return err
		}

		log.Printf("Deleted %d messages for chat %d", deletedCount, chatID)

		deleted, err := repo.DeleteChatById(ctx, chatID)
		if err != nil {
			return err
		}

		if !deleted {
			log.Printf("Warning: Chat %d was not found during deletion", chatID)
		}

		return nil
	})

	if err != nil {
		log.Printf("Failed to cleanup messages for chat %d: %v", chatID, err)
	}

	log.Printf("Successfully cleaned up all data for chat %d", chatID)
}
