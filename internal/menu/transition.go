package menu

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (mm *MenuManager) TransitionTo(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	state *MenuState,
	next State,
) {
	state.PrevState = state.State
	state.State = next

	switch next {
	case StateMain:
		mm.renderMain(ctx, b, state)
	case StateChats:
		mm.renderChats(ctx, b, state)
	case StateChatActions:
		mm.renderChatActions(ctx, b, state)
	case StateDuplicateMessage:
		mm.renderDuplicateAction(ctx, b, state)
	case StateSettings:
		mm.renderSettings(ctx, b, state)
	case StateFeatures:
		mm.renderFeatures(ctx, b, state)
	case StateFeatureEdit:
		mm.renderFeatureEdit(ctx, b, state)
	case StateCommands:
		mm.renderCommands(ctx, b, state)
	case StateCommandEdit:
		mm.renderCommandEdit(ctx, b, state)
	}
}
