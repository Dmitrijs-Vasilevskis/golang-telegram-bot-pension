package menu

import (
	"context"
	"fmt"
	"sort"

	"github.com/go-telegram/bot/models"
)

func (mm *MenuManager) mainKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "⚙️ Select chat", CallbackData: "nav:chats"},
			},
		},
	}
}

func (mm *MenuManager) actionKeyboard(feature string) *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "✅ Enable", CallbackData: fmt.Sprintf("features:command:%s:enable", feature)},
				{Text: "❌ Disable", CallbackData: fmt.Sprintf("features:command:%s:disable", feature)},
			},
			{
				{Text: "🔙 Return", CallbackData: "nav:commands"},
				{Text: "🏠 Main Menu", CallbackData: "nav:main"},
			},
		},
	}
}

func (mm *MenuManager) buildFeatureEditKeyboard(feature string) *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "✅ Enable", CallbackData: fmt.Sprintf("feature:%s:enable", feature)},
				{Text: "❌ Disable", CallbackData: fmt.Sprintf("feature:%s:disable", feature)},
			},
			{
				{Text: "🔙 Return", CallbackData: "nav:features"},
				{Text: "🏠 Main Menu", CallbackData: "nav:main"},
			},
		},
	}
}

func (mm *MenuManager) buildSettingKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "💾 Features", CallbackData: "nav:features"},
				{Text: "🎮 Commands", CallbackData: "nav:commands"},
			},
			{
				{Text: "🔙 Return", CallbackData: "nav:chat_selected"},
			},
		},
	}
}

func (mm *MenuManager) buildChatActionsBeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "⚙️ Chat configurations", CallbackData: "nav:settings"},
				{Text: "Duplicate message", CallbackData: "nav:action:duplicate_dm"},
			},
			{
				{Text: "🔙 Return", CallbackData: "nav:chats"},
			},
		},
	}
}

func (mm *MenuManager) buildFeaturesKeyboard() *models.InlineKeyboardMarkup {
	features := mm.commandManager.GetFeatureCommands()

	var buttons [][]models.InlineKeyboardButton

	for _, feature := range features {
		buttons = append(buttons, []models.InlineKeyboardButton{
			{Text: fmt.Sprintf("%s %s", feature.Icon, feature.Name), CallbackData: fmt.Sprintf("nav:features:select:%s", feature.Key)},
		})
	}

	buttons = append(buttons, []models.InlineKeyboardButton{
		{Text: "🔄 Change Chat", CallbackData: "nav:chats"},
		{Text: "🔙 Return to Settings", CallbackData: "nav:settings"},
	})

	return &models.InlineKeyboardMarkup{InlineKeyboard: buttons}
}

func (mm *MenuManager) buildCommandsKeyboard(ctx context.Context, state *MenuState) *models.InlineKeyboardMarkup {
	commands := mm.commandManager.GetRegularCommands()
	commandStates, _ := mm.configService.GetCommands(ctx, state.ChatID)

	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name < commands[j].Name
	})

	var buttons [][]models.InlineKeyboardButton

	buttons = append(buttons, []models.InlineKeyboardButton{
		{Text: "✅ Enable All", CallbackData: "nav:commands:enable_all"},
		{Text: "❌ Disable All", CallbackData: "nav:commands:disable_all"},
	})

	for _, cmd := range commands {
		enabled := commandStates[cmd.Key]

		icon := "❌"
		if enabled {
			icon = "✅"
		}

		buttons = append(buttons, []models.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("%s %s", icon, cmd.Name),
				CallbackData: fmt.Sprintf("nav:commands:select:%s", cmd.Key),
			},
		})
	}

	buttons = append(buttons, []models.InlineKeyboardButton{
		{Text: "🔙 Return to Settings", CallbackData: "nav:back:settings"},
	})

	return &models.InlineKeyboardMarkup{InlineKeyboard: buttons}
}

func (mm *MenuManager) buildChatsKeyboard(ctx context.Context, state MenuState) *models.InlineKeyboardMarkup {
	chats, err := mm.chatService.GetUserChats(ctx, state.UserID)
	if err != nil {
		fmt.Printf("Failed to get user chats: %v\n", err)

	}

	var buttons [][]models.InlineKeyboardButton

	for _, chat := range chats {
		buttons = append(buttons, []models.InlineKeyboardButton{
			{
				Text:         chat.Title,
				CallbackData: fmt.Sprintf("chat:select_chat:%d", chat.ID),
			},
		})
	}

	buttons = append(buttons, []models.InlineKeyboardButton{
		{Text: "🏠 Main Menu", CallbackData: "nav:back:main"},
	})

	return &models.InlineKeyboardMarkup{InlineKeyboard: buttons}
}

func (mm *MenuManager) buildDuplicateActionKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "🔙 Return", CallbackData: "nav:chat_selected"},
				{Text: "🏠 Main Menu", CallbackData: "nav:main"},
			},
		},
	}
}
