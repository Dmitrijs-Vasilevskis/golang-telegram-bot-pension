package menu

import (
	"context"
	"sync"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/app"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
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
	CurrentNode string
	Data        map[string]interface{}
}

type MenuManager struct {
	nodes     map[string]*MenuNode
	userState map[int64]*MenuState
	mu        sync.RWMutex
	app       *app.App
}

type FeatureHandler interface {
	Handle(ctx context.Context, b *bot.Bot, update *models.Update, parts []string)
}
