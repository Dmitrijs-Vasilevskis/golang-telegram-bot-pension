package handlers

import (
	"context"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/helpers"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/service"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func TikTok(ctx context.Context, b *bot.Bot, update *models.Update) {
	url := update.Message.Text

	if helpers.IsShortUrl(url) {
		resolved, err := service.ResolveShortenedUrl(url)
		if err != nil {
			return
		}

		url = resolved
	}

	url, err := helpers.ConverTikTokToVxUrl(url)
	if err != nil {
		return
	}

	utils.Reply(ctx, b, update, url)
}
