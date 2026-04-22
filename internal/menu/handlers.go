package menu

import (
	"context"
	"strings"

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

	action := callback.Data

	switch state.State {
	case StateMain:
		mm.handleMain(ctx, b, update, state, action)
	case StateChats:
		mm.handleChats(ctx, b, update, state, action)
	case StateSettings:
		mm.handleSettings(ctx, b, update, state, action)
	case StateFeatures:
		mm.handleFeatures(ctx, b, update, state, action)
	case StateFeatureEdit:
		mm.handleFeatureEdit(ctx, b, update, state, action)
	case StateCommands:
		mm.handleCommands(ctx, b, update, state, action)
	case StateCommandEdit:
		mm.handleCommandEdit(ctx, b, update, state, action)
	default:
		state.State = StateMain
		mm.renderMain(ctx, b, state)
	}
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

func (mm *MenuManager) handleMain(ctx context.Context, b *bot.Bot, update *models.Update, state *MenuState, action string) {
	callback := update.CallbackQuery
	parts := strings.Split(action, ":")

	if len(parts) < 1 {
		return
	}

	switch parts[0] {
	case "nav":
		switch parts[1] {
		case "chats":
			state.State = StateChats

			mm.renderChats(ctx, b, state)
			mm.answerCallback(ctx, b, callback.ID)
			return
		}
	}

	mm.renderMain(ctx, b, state)
	mm.answerCallback(ctx, b, callback.ID)
}

func (mm *MenuManager) handleChats(ctx context.Context, b *bot.Bot, update *models.Update, state *MenuState, action string) {
	callback := update.CallbackQuery
	parts := strings.Split(action, ":")

	if len(parts) < 1 {
		return
	}

	switch parts[0] {
	case "chat":
		if len(parts) < 2 {
			return
		}

		switch parts[1] {
		case "select_chat":
			chatID := helpers.ParseChatID(parts[2])

			if chatID == 0 {
				break
			}

			state.ChatID = chatID
			state.State = StateSettings

			mm.renderSettings(ctx, b, state)
			mm.answerCallback(ctx, b, callback.ID)
			return
		}
	case "nav":
		if len(parts) < 2 {
			return
		}

		switch parts[2] {
		case "main":
			state.State = StateMain

			mm.renderMain(ctx, b, state)
			mm.answerCallback(ctx, b, callback.ID)
			return
		}
	}

	mm.renderChats(ctx, b, state)
	mm.answerCallback(ctx, b, callback.ID)
}

func (mm *MenuManager) handleSettings(ctx context.Context, b *bot.Bot, update *models.Update, state *MenuState, action string) {
	callback := update.CallbackQuery
	parts := strings.Split(action, ":")

	if len(parts) < 1 {
		return
	}

	switch parts[1] {
	case "features":
		state.State = StateFeatures

		mm.renderFeatures(ctx, b, state)
		mm.answerCallback(ctx, b, callback.ID)
		return
	case "commands":
		state.State = StateCommands

		mm.renderCommands(ctx, b, state)
		mm.answerCallback(ctx, b, callback.ID)
		return
	case "chats":
		state.State = StateChats

		mm.renderChats(ctx, b, state)
		mm.answerCallback(ctx, b, callback.ID)
		return
	}

	mm.renderSettings(ctx, b, state)
	mm.answerCallback(ctx, b, callback.ID)
}

func (mm *MenuManager) handleFeatures(ctx context.Context, b *bot.Bot, update *models.Update, state *MenuState, action string) {
	callback := update.CallbackQuery
	parts := strings.Split(action, ":")

	if len(parts) < 1 {
		return
	}

	switch parts[1] {
	case "features":

		switch parts[2] {
		case "select":
			feature := parts[3]
			state.State = StateFeatureEdit
			state.Data["feature"] = feature

			if _, ok := mm.commandManager.GetByKey(feature); !ok {
				break
			}

			mm.renderFeatureEdit(ctx, b, state)
			mm.answerCallback(ctx, b, callback.ID)
			return
		}

		mm.renderFeatures(ctx, b, state)
		mm.answerCallback(ctx, b, callback.ID)
		return
	case "chats":
		state.State = StateChats

		mm.renderChats(ctx, b, state)
		mm.answerCallback(ctx, b, callback.ID)
		return

	case "settings":
		state.State = StateSettings

		mm.renderSettings(ctx, b, state)
		mm.answerCallback(ctx, b, callback.ID)
		return
	}

	mm.renderFeatures(ctx, b, state)
	mm.answerCallback(ctx, b, callback.ID)
}

func (mm *MenuManager) handleFeatureEdit(ctx context.Context, b *bot.Bot, update *models.Update, state *MenuState, action string) {
	callback := update.CallbackQuery
	parts := strings.Split(action, ":")

	if len(parts) < 1 {
		return
	}

	switch parts[0] {
	case "feature":
		if len(parts) < 2 {
			break
		}
		cmdName := parts[1]
		actionType := parts[2]

		switch actionType {
		case "enable":
			_ = mm.chatService.SetFeature(ctx, state.ChatID, state.UserID, cmdName, true)
		case "disable":
			_ = mm.chatService.SetFeature(ctx, state.ChatID, state.UserID, cmdName, false)
		}

		state.Data["feature"] = cmdName

		mm.renderFeatureEdit(ctx, b, state)
		mm.answerCallback(ctx, b, callback.ID)
		return
	case "nav":
		switch parts[1] {
		case "features":
			state.State = StateFeatures

			mm.renderFeatures(ctx, b, state)
			mm.answerCallback(ctx, b, callback.ID)
			return
		case "main":
			state.State = StateMain

			mm.renderMain(ctx, b, state)
			mm.answerCallback(ctx, b, callback.ID)
			return
		}
	}

	mm.renderFeatureEdit(ctx, b, state)
	mm.answerCallback(ctx, b, callback.ID)
}

func (mm *MenuManager) handleCommands(ctx context.Context, b *bot.Bot, update *models.Update, state *MenuState, action string) {
	callback := update.CallbackQuery
	parts := strings.Split(action, ":")

	if len(parts) < 1 {
		return
	}

	switch parts[0] {
	case "nav":
		switch parts[1] {
		case "commands":
			if len(parts) < 2 {
				break
			}
			switch parts[2] {
			case "enable_all":
				_ = mm.configService.SetAllCommands(ctx, state.ChatID, state.UserID, true)
			case "disable_all":
				_ = mm.configService.SetAllCommands(ctx, state.ChatID, state.UserID, false)
			case "select":
				if len(parts) < 3 {
					break
				}

				cmd := parts[3]

				if _, ok := mm.commandManager.GetByKey(cmd); !ok {
					break
				}

				state.State = StateCommandEdit
				state.Data["command"] = cmd

				mm.renderCommandEdit(ctx, b, state)
				mm.answerCallback(ctx, b, callback.ID)
				return
			}
		case "back":
			if len(parts) < 2 {
				break
			}
			switch parts[2] {
			case "settings":
				state.State = StateSettings

				mm.renderSettings(ctx, b, state)
				mm.answerCallback(ctx, b, callback.ID)
				return
			}
		}
	}

	mm.renderCommands(ctx, b, state)
	mm.answerCallback(ctx, b, callback.ID)
}

func (mm *MenuManager) handleCommandEdit(ctx context.Context, b *bot.Bot, update *models.Update, state *MenuState, action string) {
	callback := update.CallbackQuery
	parts := strings.Split(action, ":")

	if len(parts) < 1 {
		return
	}

	switch parts[0] {
	case "features":
		if parts[1] != "command" {
			break
		}

		cmd := parts[2]
		actionType := parts[3]

		switch actionType {
		case "enable":
			_ = mm.configService.SetCommand(ctx, state.ChatID, state.UserID, cmd, true)
		case "disable":
			_ = mm.configService.SetCommand(ctx, state.ChatID, state.UserID, cmd, false)
		}

		state.Data["command"] = cmd

		mm.renderCommandEdit(ctx, b, state)
		mm.answerCallback(ctx, b, callback.ID)
		return
	case "nav":
		switch parts[1] {
		case "commands":
			state.State = StateCommands

			mm.renderCommands(ctx, b, state)
			mm.answerCallback(ctx, b, callback.ID)
			return
		case "main":
			state.State = StateMain

			mm.renderMain(ctx, b, state)
			mm.answerCallback(ctx, b, callback.ID)
			return
		}
	}

	mm.renderCommandEdit(ctx, b, state)
	mm.answerCallback(ctx, b, callback.ID)
}
