package ask

import (
	"context"
	"fmt"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/helpers"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/service"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Handler(gemini *service.GeminiService) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		message := update.Message

		mediaData, err := helpers.ProcessMedia(update)
		if err != nil {
			if helpers.HasMedia(message) ||
				(message.ReplyToMessage != nil && helpers.HasMedia(message.ReplyToMessage)) {
				utils.Reply(ctx, b, update, "Неподдерживаемый тип медиа.")
				return
			}
			// fallback to "no media"
			mediaData = nil
		}

		prompt, err := service.BuildPromptInput(message, mediaData != nil)
		if err != nil {
			utils.Reply(ctx, b, update, "Введите запрос.")
			return
		}

		if mediaData != nil {
			token := b.Token()

			fileLink, err := service.GetDownloadLink(b, helpers.GetFileUrl(mediaData.FileId, token))
			if err != nil {
				fmt.Printf("[WARN] Failed to get download link for file %s: %v\n", mediaData.FileId, err)
				utils.Reply(ctx, b, update, "Не удалось получить ссылку для скачивания. Попробуйте позже.")
				return
			}

			response, err := gemini.GenResponseWithMediaPreset(ctx, prompt, fileLink, service.PromptTypeCustom)
			if err != nil {
				fmt.Println("[WARN] Gemini API failed:", err)
				utils.Reply(ctx, b, update, "Не удалось обработать изображение. Возможно, формат не поддерживается или сервис временно недоступен.")
				return
			}

			utils.Reply(ctx, b, update, response)
		} else {
			response, err := gemini.GenResponseWithPreset(ctx, prompt, service.PromptTypeCustom)
			if err != nil {
				fmt.Println("[ERROR] Gemini API failed:", err)
				utils.Reply(ctx, b, update, "⚠️ Не удалось обработать запрос. Попробуйте позже.")
				return
			}

			utils.Reply(ctx, b, update, response)
		}
	}
}
