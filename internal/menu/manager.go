package menu

import (
	"context"
	"fmt"
	"strings"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/app"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func NewMenuManager(app *app.App) *MenuManager {
	mm := &MenuManager{
		nodes:     make(map[string]*MenuNode),
		userState: make(map[int64]*MenuState),
		app:       app,
	}

	mm.initNodes()

	return mm
}

func (mm *MenuManager) initNodes() {
	mm.nodes["main"] = &MenuNode{
		ID:       "main",
		Text:     "*Main Menu*",
		ParentID: "",
		Keyboard: mm.mainKeyboard(),
	}

	mm.nodes["settings"] = &MenuNode{
		ID:       "settings",
		Text:     "*Bot Settings*",
		ParentID: "main",
		Keyboard: mm.settingKeyboard(),
	}

	mm.nodes["chats"] = &MenuNode{
		ID:       "chats",
		Text:     "Select a chat:",
		ParentID: "main",
	}

	mm.nodes["features"] = &MenuNode{
		ID:       "features",
		Text:     "Select a feature to configure",
		ParentID: "settings",
		Keyboard: mm.featureKeyboard(),
	}

	mm.nodes["commands"] = &MenuNode{
		ID:       "commands",
		Text:     "Command Management",
		ParentID: "features",
		Keyboard: mm.commandsKeyboard(),
	}

	mm.nodes["summary"] = &MenuNode{
		ID:       "summary",
		Text:     "Summary",
		ParentID: "features",
		Keyboard: mm.featureToggleKeyboard("summary"),
	}

	mm.nodes["duplicate_dm"] = &MenuNode{
		ID:       "duplicate_dm",
		Text:     "Duplicate DM",
		ParentID: "features",
		Keyboard: mm.featureToggleKeyboard("duplicate_dm"),
	}
}

func (mm *MenuManager) GetNode(nodeID string) (*MenuNode, bool) {
	node, exists := mm.nodes[nodeID]
	return node, exists
}

func (mm *MenuManager) ShowMenu(ctx context.Context, b *bot.Bot, chatID int, userID int, nodeID string) error {
	node, exists := mm.GetNode(nodeID)
	if !exists {
		return fmt.Errorf("node %s not found", nodeID)
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        node.Text,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: node.Keyboard,
	})

	return err
}

func (mm *MenuManager) updateMenu(ctx context.Context, b *bot.Bot, chatID int, messageID int64, nodeID string) error {
	node, exists := mm.GetNode(nodeID)
	if !exists {
		return fmt.Errorf("node %s not found", nodeID)
	}

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatID,
		MessageID:   int(messageID),
		Text:        node.Text,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: node.Keyboard,
	})

	return err
}

func (mm *MenuManager) HandleCallback(ctx context.Context, b *bot.Bot, update *models.Update, app *app.App) {
	callback := update.CallbackQuery
	data := callback.Data

	parts := strings.Split(data, ":")
	if len(parts) < 2 {
		return
	}

	logger.DebugJson(callback)

	switch parts[0] {
	case "menu":
		mm.handleMenuNavigation(ctx, b, update, parts)
	case "features":

		mm.handleFeatureAction(ctx, b, update, parts)
	default:
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callback.ID,
		})
	}
}

func (mm *MenuManager) setUserChat(userID int64, chatID int64) {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	state, ok := mm.userState[userID]
	if !ok {
		state = &MenuState{
			UserID: userID,
			Data:   make(map[string]interface{}),
		}
		mm.userState[userID] = state
	}

	state.ChatID = chatID
}

func (mm *MenuManager) getUserChat(userID int64) int64 {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	if state, ok := mm.userState[userID]; ok {
		return state.ChatID
	}

	return 0
}
