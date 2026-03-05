package router

import (
	"context"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type MessageHandler func(ctx context.Context, b *bot.Bot, update *models.Update)

type Router struct {
	handlers map[string]MessageHandler
}

func NewRouter() *Router {
	return &Router{
		handlers: make(map[string]MessageHandler),
	}
}

func (r *Router) Register(keyword string, handler MessageHandler) {
	r.handlers[keyword] = handler
}

func (r *Router) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := update.Message.Text

	for keyword, handler := range r.handlers {
		if strings.Contains(text, keyword) {
			handler(ctx, b, update)

			return
		}
	}
}
