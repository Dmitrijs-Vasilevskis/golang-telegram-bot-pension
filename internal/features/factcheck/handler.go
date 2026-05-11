package factcheck

import (
	"context"
	"fmt"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/service"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Handler(gemini *service.GeminiService) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		message := update.Message

		prompt, err := service.BuildPromptInput(message, false)
		if err != nil {
			fmt.Println("[WARN] No claim provided by user")
			utils.Reply(ctx, b, update, "Пожалуйста, укажите утверждение после /factcheck или выберите сообщение для ответа")
			return
		}

		fmt.Println("[INFO] User claim:", prompt)
		check, err := gemini.GenResponseWithPreset(ctx, prompt, service.PromptTypeFactCheck)
		if err != nil {
			fmt.Println("[ERROR] Gemini API failed:", err)
			utils.Reply(ctx, b, update, "⚠️ Не удалось проверить факт. Попробуйте позже.")
			return
		}
		fmt.Println("[INFO] Gemini response received")

		utils.Reply(ctx, b, update, check)
	}
}
