package menu

import (
	"context"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/helpers"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (mm *MenuManager) HandleCallback(ctx context.Context, b *bot.Bot, update *models.Update) {
	callback := update.CallbackQuery
	userID := callback.From.ID

	state := mm.getOrCreateState(userID)

	if !mm.trySetLoading(userID) {
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callback.ID,
			Text:            "⏳ Please wait...",
		})
		return
	}

	defer mm.setLoading(userID, false)

	mm.router.Dispatch(ctx, b, update, state)
	mm.answerCallback(ctx, b, callback.ID)
}

func (mm *MenuManager) HandleStart(ctx context.Context, b *bot.Bot, update *models.Update) {
	user := update.Message.From
	chat := update.Message.Chat

	state := mm.resetState(user.ID)
	state.DMChannelID = chat.ID

	msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chat.ID,
		Text:        "*Main menu*",
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: mm.mainKeyboard(),
	})

	if err != nil {
		return
	}

	state.MessageId = msg.ID
	state.State = StateMain
}

func (mm *MenuManager) HandleMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	userID := update.Message.From.ID
	state := mm.getOrCreateState(userID)

	switch update.Message.Text {
	case "/cancel":
		mm.handleCancel(ctx, b, update, state)
	case "/reload":
		mm.handleReload(ctx, b, update, state)
	}

	mm.messageRouter.Dispatch(ctx, b, update, state)
}

func (mm *MenuManager) handleReload(ctx context.Context, b *bot.Bot, update *models.Update, state *MenuState) {
	chatID := update.Message.Chat.ID

	state.DMChannelID = chatID

	msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "🔄 Reloading menu...",
	})

	if err != nil {
		return
	}

	state.MessageId = msg.ID

	mm.TransitionTo(ctx, b, update, state, state.State)
}

func (mm *MenuManager) handleCancel(ctx context.Context, b *bot.Bot, update *models.Update, state *MenuState) {
	if state.PrevState != "" {
		mm.TransitionTo(ctx, b, update, state, state.PrevState)
	}
}

func (mm *MenuManager) handleMain(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	state *MenuState,
	parts []string,
) {
	mm.TransitionTo(ctx, b, update, state, StateMain)

}

func (mm *MenuManager) handleChats(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	state *MenuState,
	parts []string,
) {
	mm.TransitionTo(ctx, b, update, state, StateChats)
}

func (mm *MenuManager) handleFeatures(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	state *MenuState,
	parts []string,
) {
	mm.TransitionTo(ctx, b, update, state, StateFeatures)
}

func (mm *MenuManager) handleCommands(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	state *MenuState,
	parts []string,
) {
	mm.TransitionTo(ctx, b, update, state, StateCommands)
}

func (mm *MenuManager) handleDuplicateMessage(ctx context.Context, b *bot.Bot, update *models.Update, state *MenuState) {
	msg := update.Message

	if msg == nil {
		return
	}

	if state.ChatID == 0 {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: state.DMChannelID,
			Text:   "❌ The chat is not selected",
		})
		return
	}

	enabled, _ := mm.chatService.IsDuplicateDMEnabled(ctx, state.ChatID)

	if !enabled {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: state.DMChannelID,
			Text:   "❌ The duplicate feature is not enabled",
		})
		return
	}

	_, err := b.CopyMessage(ctx, &bot.CopyMessageParams{
		ChatID:     state.ChatID,
		FromChatID: msg.Chat.ID,
		MessageID:  msg.ID,
	})

	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: state.DMChannelID,
			Text:   "❌ Failed to duplicate message",
		})
		return
	}
}

func (mm *MenuManager) handleDuplicateMode(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	state *MenuState,
	parts []string,
) {
	mm.TransitionTo(ctx, b, update, state, StateDuplicateMessage)
}

func (mm *MenuManager) handleSelectChat(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	state *MenuState,
	parts []string,
) {
	if len(parts) < 4 {
		return
	}

	chatID := helpers.ParseChatID(parts[3])

	if chatID == 0 {
		return
	}

	state.ChatID = chatID

	mm.TransitionTo(ctx, b, update, state, StateChatActions)
}

func (mm *MenuManager) handleSelectedChat(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	state *MenuState,
	parts []string,
) {
	if len(parts) < 3 {
		return
	}

	chatID := state.ChatID

	if chatID == 0 {
		return
	}

	mm.TransitionTo(ctx, b, update, state, StateChatActions)
}

func (mm *MenuManager) handleNavSettings(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	state *MenuState,
	parts []string,
) {
	if len(parts) < 2 {
		return
	}

	mm.TransitionTo(ctx, b, update, state, StateSettings)
}

func (mm *MenuManager) handleSelectCommand(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	state *MenuState,
	parts []string,
) {
	if len(parts) < 4 {
		return
	}

	cmdKey := parts[3]

	if _, ok := mm.commandManager.GetByKey(cmdKey); !ok {
		return
	}

	state.Data["command"] = cmdKey

	mm.TransitionTo(ctx, b, update, state, StateCommandEdit)
}

func (mm *MenuManager) handleAllCommands(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	state *MenuState,
	parts []string,
) {
	if len(parts) < 4 {
		return
	}

	actionType := parts[3]

	switch actionType {
	case "enable":
		_ = mm.configService.SetAllCommands(ctx, state.ChatID, state.UserID, true)
	case "disable":
		_ = mm.configService.SetAllCommands(ctx, state.ChatID, state.UserID, false)
	default:
		return
	}

	mm.TransitionTo(ctx, b, update, state, StateCommands)
}

func (mm *MenuManager) handleToggleCommand(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	state *MenuState,
	parts []string,
) {
	if len(parts) < 4 {
		return
	}

	cmdKey := parts[2]
	actionType := parts[3]

	switch actionType {
	case "enable":
		_ = mm.configService.SetCommand(ctx, state.ChatID, state.UserID, cmdKey, true)
	case "disable":
		_ = mm.configService.SetCommand(ctx, state.ChatID, state.UserID, cmdKey, false)
	default:
		return
	}

	state.Data["command"] = cmdKey

	mm.TransitionTo(ctx, b, update, state, StateCommandEdit)
}

func (mm *MenuManager) handleSelectFeature(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	state *MenuState,
	parts []string,
) {
	if len(parts) < 4 {
		return
	}

	featureKey := parts[3]

	if _, ok := mm.commandManager.GetByKey(featureKey); !ok {
		return
	}

	state.Data["feature"] = featureKey

	mm.TransitionTo(ctx, b, update, state, StateFeatureEdit)
}

func (mm *MenuManager) handleToggleFeature(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	state *MenuState,
	parts []string,
) {
	if len(parts) < 4 {
		return
	}

	featureKey := parts[2]
	actionType := parts[3]

	switch actionType {
	case "enable":
		_ = mm.chatService.SetFeature(ctx, state.ChatID, state.UserID, featureKey, true)
	case "disable":
		_ = mm.chatService.SetFeature(ctx, state.ChatID, state.UserID, featureKey, false)
	default:
		return
	}

	state.Data["feature"] = featureKey

	mm.TransitionTo(ctx, b, update, state, StateFeatureEdit)
}
