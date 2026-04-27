package menu

import (
	"context"
	"sync"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/command"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/service"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type State string

const (
	StateMain             State = "main"
	StateChats            State = "chats"
	StateChatActions      State = "chat_actions"
	StateSettings         State = "settings"
	StateFeatures         State = "features"
	StateFeatureEdit      State = "feature_edit"
	StateCommands         State = "commands"
	StateCommandEdit      State = "command_edit"
	StateDuplicateMessage State = "duplicate_dm"
)

type MenuNode struct {
	ID       string
	Text     string
	ParentID string
	Keyboard *models.InlineKeyboardMarkup
	Handler  func(ctx context.Context, b *bot.Bot, update *models.Update)
}

type MenuState struct {
	UserID      int64
	ChatID      int64
	DMChannelID int64
	State       State
	MessageId   int
	Loading     bool
	LastAction  string
	CurrentNode string
	Data        map[string]interface{}
}

type MenuManager struct {
	nodes          map[string]*MenuNode
	userState      map[int64]*MenuState
	commandManager command.CommandManager
	mu             sync.RWMutex
	configService  service.ConfigService
	chatService    service.ChatService
}

type FeatureHandler interface {
	Handle(ctx context.Context, b *bot.Bot, update *models.Update, parts []string)
}
