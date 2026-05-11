package summary

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/helpers"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/repository"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/service"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Handler(repo *repository.Repository, gemini *service.GeminiService) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		chatID := update.Message.Chat.ID
		userText := update.Message.Text

		limit := 200
		const maxLimit = 400

		if args, hasArgs := helpers.GetCommandArgs(userText); hasArgs {
			parsedLimit, err := strconv.Atoi(args)
			if err != nil {
				utils.Reply(ctx, b, update, "Invalid argument. Please provide a number ")
				return
			}

			if parsedLimit < 1 || parsedLimit > maxLimit {
				utils.Reply(ctx, b, update, fmt.Sprintf("Limit must be between 1 and %d", maxLimit))
				return
			}

			limit = parsedLimit
		}

		messages, err := repo.GetLastMessages(ctx, chatID, limit)
		if err != nil {
			fmt.Println("[WARN] Failed to get messages:", err)
			utils.Reply(ctx, b, update, "Failed to retrieve messages.")
			return
		}

		processedMessages, err := helpers.ProcessSummary(messages)
		if err != nil {
			fmt.Println("[WARN] Failed to process messages:", err)
			utils.Reply(ctx, b, update, "Failed to process messages.")
			return
		}

		summaryJSON, err := json.Marshal(processedMessages)
		if err != nil {
			fmt.Println("[WARN] Failed to marshal summary:", err)
			utils.Reply(ctx, b, update, "Failed to prepare summary.")
			return
		}

		ans, err := gemini.GenResponseWithPreset(ctx, string(summaryJSON), service.PromptTypeSummary)
		if err != nil {
			fmt.Println("[WARN] Failed to generate summary:", err)
			utils.Reply(ctx, b, update, "Failed to generate summary with AI.")
			return
		}

		utils.Reply(ctx, b, update, ans)
	}
}
