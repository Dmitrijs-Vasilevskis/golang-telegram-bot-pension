package handlers

import (
	"context"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Clown(ctx context.Context, b *bot.Bot, update *models.Update) {
	utils.Emote(ctx, b, update, false, "🤡")
}
