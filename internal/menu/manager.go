package menu

import (
	"context"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/app"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/command"
	"github.com/go-telegram/bot"
)

func NewMenuManager(app *app.App) *MenuManager {
	mm := &MenuManager{
		userState:      make(map[int64]*MenuState),
		router:         NewCallbackRouter(),
		messageRouter:  NewMessagerRouter(),
		commandManager: *command.New(),
		configService:  app.Services.Config,
		chatService:    app.Services.Chat,
	}

	mm.registerCallbackRoutes()
	mm.registerMessageRoutes()

	return mm
}

func (mm *MenuManager) getOrCreateState(userId int64) *MenuState {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	if state, ok := mm.userState[userId]; ok {
		return state
	}

	state := &MenuState{
		UserID: userId,
		Data:   make(map[string]interface{}),
	}

	mm.userState[userId] = state

	return state
}

func (mm *MenuManager) resetState(userID int64) *MenuState {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	state := &MenuState{
		UserID: userID,
		Data:   make(map[string]interface{}),
	}

	mm.userState[userID] = state

	return state
}

func (s *MenuState) SetCommand(cmd string) {
	s.Data["command"] = cmd
}

func (s *MenuState) GetCommand() string {
	if v, ok := s.Data["command"].(string); ok {
		return v
	}
	return ""
}

func (mm *MenuManager) answerCallback(ctx context.Context, b *bot.Bot, callbackID string) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackID,
	})
}

func (mm *MenuManager) setLoading(userID int64, loading bool) {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	if state, ok := mm.userState[userID]; ok {
		state.Loading = loading
	}
}

func (mm *MenuManager) trySetLoading(userID int64) bool {
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

	if state.Loading {
		return false
	}

	state.Loading = true
	return true
}
