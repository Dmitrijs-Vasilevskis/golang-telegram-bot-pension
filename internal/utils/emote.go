package utils

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Emote(ctx context.Context, b *bot.Bot, update *models.Update, big bool, emoji string) {
	reaction := []models.ReactionType{
		{
			Type: models.ReactionTypeTypeEmoji,
			ReactionTypeEmoji: &models.ReactionTypeEmoji{
				Emoji: emoji,
			},
		},
	}

	_, err := b.SetMessageReaction(ctx, &bot.SetMessageReactionParams{
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.ID,
		Reaction:  reaction,
		IsBig:     &big,
	})

	if err != nil {
		log.Println("failed to set reaction", err)
	}
}
