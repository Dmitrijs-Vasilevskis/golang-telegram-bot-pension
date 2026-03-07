package handlers

import (
	"context"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/helpers"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Reels(ctx context.Context, b *bot.Bot, update *models.Update) {
	reelId, isValid := helpers.ExtractReelId(update.Message.Text)

	if isValid {
		url, _ := helpers.BuildKkInstagramUrl(reelId)

		utils.Reply(ctx, b, update, url)
		utils.Delete(ctx, b, update)
	}
}
