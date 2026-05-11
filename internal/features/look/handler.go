package look

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
		token := b.Token()

		userText := update.Message.Caption
		userComment, _ := helpers.GetCommandArgs(userText)

		mediaData, err := helpers.ProcessMedia(update)
		if err != nil {
			fmt.Println("[WARN] Failed to process media")
			utils.Reply(ctx, b, update, "Не удалось обработать медиафайл.")
			return
		}

		fileLink, err := service.GetDownloadLink(b, helpers.GetFileUrl(mediaData.FileId, token))
		if err != nil {
			fmt.Println("[WARN]", fmt.Sprintf("Failed to get download link for file %s: %v", mediaData.FileId, err))
			utils.Reply(ctx, b, update, "Не удалось получить ссылку для скачивания. Попробуйте позже.")
			return
		}

		ans, err := gemini.GenResponseWithMediaPreset(ctx, userComment, fileLink, service.PromptTypeAnalyze)
		if err != nil {
			fmt.Println("[WARN] Gemini API failed:", err)
			utils.Reply(ctx, b, update, "Не удалось обработать изображение. Возможно, формат не поддерживается или сервис временно недоступен.")
			return
		}

		utils.Reply(ctx, b, update, ans)
	}
}
