package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/app"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/helpers"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/service"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Summary(ctx context.Context, b *bot.Bot, update *models.Update, app *app.App) {
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

	geminiService := service.NewGeminiService()

	messages, err := app.Repository.GetLastMessages(ctx, chatID, limit)
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

	ans, err := geminiService.GenResponseWithPreset(ctx, string(summaryJSON), service.PromptTypeSummary)
	if err != nil {
		fmt.Println("[WARN] Failed to generate summary:", err)
		utils.Reply(ctx, b, update, "Failed to generate summary with AI.")
		return
	}

	utils.Reply(ctx, b, update, ans)
}
