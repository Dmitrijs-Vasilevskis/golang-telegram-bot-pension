package duplicate

import (
	"context"
	"fmt"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type MenuHandler struct {
	feature *Feature
}

func NewMenuHandler(feature *Feature) *MenuHandler {
	return &MenuHandler{feature: feature}
}

func (h *MenuHandler) Handle(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	chatID int64,
	dmChannedID int64,
) {
	msg := update.Message

	if msg == nil {
		return
	}

	if chatID == 0 {
		utils.Reply(ctx, b, update, "❌ No chat selected")
	}

	err := h.feature.CopyMessage(ctx, b, chatID, msg.Chat.ID, msg.ID, dmChannedID)
	if err != nil {
		fmt.Printf("[ERROR] Failed to copy message: %v\n", err)
		utils.Reply(ctx, b, update, "❌ Failed to copy message")

		return
	}
}
