package menu

import (
	"context"
	"fmt"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/helpers"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/service"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (mm *MenuManager) handleMenuNavigation(ctx context.Context, b *bot.Bot, update *models.Update, parts []string) {
	callback := update.CallbackQuery
	chatID := callback.Message.Message.Chat.ID
	userID := callback.From.ID
	messageId := callback.Message.Message.ID

	if len(parts) < 2 {
		return
	}

	action := parts[1]

	switch action {
	case "back":
		if len(parts) < 3 {
			return
		}
		targetNode := parts[2]

		if err := mm.updateMenu(ctx, b, int(chatID), int64(messageId), targetNode); err != nil {
			fmt.Printf("Failed to update menu: %v\n", err)
		}

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callback.ID,
		})
	case "select_chat":
		chatID := helpers.ParseChatID(parts[2])
		mm.setUserChat(userID, chatID)

		if err := mm.updateMenu(ctx, b, int(callback.Message.Message.Chat.ID), int64(messageId), "settings"); err != nil {
			fmt.Printf("Failed to update menu: %v\n", err)
		}

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callback.ID,
		})
	case "settings":
		mm.showChatList(ctx, b, chatID, int64(messageId), userID)

	default:
		if err := mm.updateMenu(ctx, b, int(chatID), int64(messageId), action); err != nil {
			fmt.Printf("Failed to update menu: %v\n", err)
		}

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callback.ID,
		})
	}
}

func (mm *MenuManager) handleFeatureAction(ctx context.Context, b *bot.Bot, update *models.Update, parts []string) {
	callback := update.CallbackQuery
	userID := callback.From.ID
	selectedChatID := mm.getUserChat(userID)

	// if selectedChatID == 0 {
	// 	return
	// }

	if len(parts) < 3 {
		return
	}

	dmChatID := callback.Message.Message.Chat.ID
	messageID := callback.Message.Message.ID

	cfg := service.NewConfigService(mm.app.Repository)

	feature := parts[1]
	action := parts[2]

	println("handleFeatureAction parts[1]", parts[1])
	println("handleFeatureAction parts[2]", parts[2])

	switch feature {
	case "commands":
		switch action {
		case "enable_all":
			_ = cfg.SetAllCommands(ctx, selectedChatID, userID, true)
			_ = mm.updateMenu(ctx, b, int(dmChatID), int64(messageID), "commands")
		case "disable_all":
			_ = cfg.SetAllCommands(ctx, selectedChatID, userID, false)
			_ = mm.updateMenu(ctx, b, int(dmChatID), int64(messageID), "commands")
		default:
			// action is command name: features:commands:<cmd>
			println(">> default showCommandActionKeyboard:", action)
			mm.showCommandActionKeyboard(ctx, b, dmChatID, messageID, action)
		}

	case "command":
		// features:command:<cmd>:enable|disable
		if len(parts) < 4 {
			return
		}
		cmd := parts[2]
		cmdAction := parts[3]
		switch cmdAction {
		case "enable":
			_ = cfg.SetCommandEnabled(ctx, selectedChatID, userID, cmd, true)
		case "disable":
			_ = cfg.SetCommandEnabled(ctx, selectedChatID, userID, cmd, false)
		}
		mm.showCommandActionKeyboard(ctx, b, dmChatID, messageID, cmd)

	case "summary":
		switch action {
		case "enable":
			_ = cfg.SetSummary(ctx, selectedChatID, userID, true)
		case "disable":
			_ = cfg.SetSummary(ctx, selectedChatID, userID, false)
		}
		_ = mm.updateMenu(ctx, b, int(dmChatID), int64(messageID), "summary")

	case "duplicate_dm":
		switch action {
		case "enable":
			_ = cfg.SetDuplicateDM(ctx, selectedChatID, userID, true)
		case "disable":
			_ = cfg.SetDuplicateDM(ctx, selectedChatID, userID, false)
		}
		_ = mm.updateMenu(ctx, b, int(dmChatID), int64(messageID), "duplicate_dm")
	}

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callback.ID,
	})
}

func (mm *MenuManager) showCommandActionKeyboard(ctx context.Context, b *bot.Bot, dmChatID int64, messageID int, commandName string) {
	println(">> showCommandActionKeyboard:", commandName)
	keyboard := mm.actionKeyboard(commandName)

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      dmChatID,
		MessageID:   int(messageID),
		ParseMode:   models.ParseModeMarkdown,
		Text:        fmt.Sprintf("*Command:* `/%s`\n\nChoose status:", commandName),
		ReplyMarkup: keyboard,
	})
	if err != nil {
		fmt.Printf("Failed to show command action keyboard: %v\n", err)
	}
}

func (mm *MenuManager) showChatList(ctx context.Context, b *bot.Bot, chatID, messageID, userID int64) {
	chats, err := mm.app.Services.Chat.GetUserChats(ctx, userID)
	if err != nil {
		fmt.Printf("Failed to get user chats: %v\n", err)
		mm.updateMenu(ctx, b, int(chatID), messageID, "main")
		return
	}

	if len(chats) == 0 {
		keyboard := &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "🔙 Back", CallbackData: "menu:back:main"},
				},
			},
		}

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      chatID,
			MessageID:   int(messageID),
			Text:        "You don't have any chats where you are an administrator.\n\nAdd the bot to a group and make me an admin to configure settings.",
			ParseMode:   models.ParseModeMarkdown,
			ReplyMarkup: keyboard,
		})
		return
	}

	node, _ := mm.GetNode("chats")
	keyboard := mm.chatsKeyboard(chats)

	b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatID,
		MessageID:   int(messageID),
		Text:        node.Text,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: keyboard,
	})
}
