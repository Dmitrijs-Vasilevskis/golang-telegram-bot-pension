package menu

import (
	"fmt"

	"github.com/go-telegram/bot/models"
)

func (mm *MenuManager) mainKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "⚙️ Settings", CallbackData: "menu:settings"},
			},
		},
	}
}

func (mm *MenuManager) chatsKeyboard(chats []models.Chat) *models.InlineKeyboardMarkup {
	var buttons [][]models.InlineKeyboardButton

	for _, chat := range chats {
		buttons = append(buttons, []models.InlineKeyboardButton{
			{
				Text:         chat.Title,
				CallbackData: fmt.Sprintf("menu:select_chat:%d", chat.ID),
			},
		})
	}

	buttons = append(buttons, []models.InlineKeyboardButton{
		{Text: "🔙 Back", CallbackData: "menu:back:main"},
	})

	return &models.InlineKeyboardMarkup{InlineKeyboard: buttons}
}

func (mm *MenuManager) settingKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Features configurations", CallbackData: "menu:features"},
				{Text: "🔙 Return", CallbackData: "menu:back:main"},
			},
		},
	}
}

func (mm *MenuManager) featureKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "💾 Summary", CallbackData: "menu:summary"},
				{Text: "🔄 Duplicate DM", CallbackData: "menu:duplicate_dm"},
				{Text: "🎮 Command configurations", CallbackData: "menu:commands"},
				{Text: "🔙 Return", CallbackData: "menu:back:settings"},
			},
		},
	}
}

func (mm *MenuManager) commandsKeyboard() *models.InlineKeyboardMarkup {
	commands := []string{"ask", "summary", "factcheck", "look"}
	var buttons [][]models.InlineKeyboardButton

	buttons = append(buttons, []models.InlineKeyboardButton{
		{Text: "✅ Enable All", CallbackData: "features:commands:enable_all"},
		{Text: "❌ Disable All", CallbackData: "features:commands:disable_all"},
	})

	for _, cmd := range commands {
		buttons = append(buttons, []models.InlineKeyboardButton{
			{Text: fmt.Sprintf("/%s", cmd), CallbackData: fmt.Sprintf("features:commands:%s", cmd)},
		})
	}

	buttons = append(buttons, []models.InlineKeyboardButton{
		{Text: "🔙 Return", CallbackData: "menu:back:features"},
	})

	return &models.InlineKeyboardMarkup{InlineKeyboard: buttons}
}

func (mm *MenuManager) actionKeyboard(feature string) *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "✅ Enable", CallbackData: fmt.Sprintf("features:command:%s:enable", feature)},
				{Text: "❌ Disable", CallbackData: fmt.Sprintf("features:command:%s:disable", feature)},
			},
			{
				{Text: "🔙 Return", CallbackData: "menu:back:commands"},
				{Text: "🏠 Main Menu", CallbackData: "menu:main"},
			},
		},
	}
}

func (mm *MenuManager) featureToggleKeyboard(feature string) *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "✅ Enable", CallbackData: fmt.Sprintf("features:%s:enable", feature)},
				{Text: "❌ Disable", CallbackData: fmt.Sprintf("features:%s:disable", feature)},
			},
			{
				{Text: "🔙 Return", CallbackData: "menu:back:features"},
				{Text: "🏠 Main Menu", CallbackData: "menu:main"},
			},
		},
	}
}
