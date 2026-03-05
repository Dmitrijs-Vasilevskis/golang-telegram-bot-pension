package utils

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Sticker(ctx context.Context, b *bot.Bot, update *models.Update, stickerFileId models.InputFile) {
	b.SendSticker(ctx, &bot.SendStickerParams{
		ChatID: update.Message.Chat.ID,
		ReplyParameters: &models.ReplyParameters{
			MessageID: update.Message.ID,
		},
		Sticker: stickerFileId,
	})
}
