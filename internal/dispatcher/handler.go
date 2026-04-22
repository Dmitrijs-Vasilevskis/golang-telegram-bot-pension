package dispatcher

import (
	"context"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/app"
	messageHandler "github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/handlers/messages"
	moderHandler "github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/handlers/moderation"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/helpers"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/middleware"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/router"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func MainHandler(app *app.App, r *router.Router, mw *middleware.Middleware) bot.HandlerFunc {
	return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
		if update == nil || update.Message == nil {
			return
		}

		message := update.Message

		messageText := message.Text
		if messageText == "" && message.Caption != "" {
			messageText = message.Caption
		}

		cmd := helpers.ParseCommand(messageText)
		chatID := message.Chat.ID

		if handler, exists := commandHandlers[cmd]; exists {
			enabled, err := mw.CommandMiddleware(ctx, chatID, cmd)
			if err != nil {
				return
			}

			if !enabled {
				return
			}

			handler(ctx, bot, update, app)
			return
		}

		if helpers.IsToxic(messageText) {
			moderHandler.Clown(ctx, bot, update)
		}

		r.Handle(ctx, bot, update)

		enabled, err := mw.SummaryMiddleware(ctx, chatID)
		if err != nil {
			return
		}

		if enabled {
			messageHandler.RecordMessage(ctx, bot, update, app)
		}
	}
}
