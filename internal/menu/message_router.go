package menu

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type MessageHandler func(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	state *MenuState,
)

type MessageRouter struct {
	routes map[State]MessageHandler
}

func NewMessagerRouter() *MessageRouter {
	return &MessageRouter{
		routes: make(map[State]MessageHandler),
	}
}

func (r *MessageRouter) Register(state State, handler MessageHandler) {
	r.routes[state] = handler
}

func (r *MessageRouter) Dispatch(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	state *MenuState,
) {
	if handler, ok := r.routes[state.State]; ok {
		handler(ctx, b, update, state)
	}
}
