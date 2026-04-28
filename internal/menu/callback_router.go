package menu

import (
	"context"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type CallbackHandler func(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	state *MenuState,
	action []string,
)

type CallbackRouter struct {
	routes map[string]CallbackHandler
}

func NewCallbackRouter() *CallbackRouter {
	return &CallbackRouter{
		routes: make(map[string]CallbackHandler),
	}
}

func (r *CallbackRouter) Register(route string, handler CallbackHandler) {
	r.routes[route] = handler
}

func (r *CallbackRouter) Dispatch(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	state *MenuState,
) {
	parts := strings.Split(update.CallbackQuery.Data, ":")

	for i := len(parts); i > 0; i-- {
		route := strings.Join(parts[:i], ":")

		if handler, ok := r.routes[route]; ok {
			handler(ctx, b, update, state, parts)
			return
		}
	}
}
