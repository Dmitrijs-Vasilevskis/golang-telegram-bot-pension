package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/app"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/command"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/database"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/dispatcher"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/features/ask"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/features/duplicate"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/features/factcheck"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/features/look"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/features/status"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/features/summary"
	chatHandler "github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/handlers/chat"
	embedHandler "github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/handlers/embed"
	messageHandler "github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/handlers/messages"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/logger"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/menu"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/middleware"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/repository"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/router"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/service"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	dbUrl := os.Getenv("DATABASE_URL")

	if token == "" {
		log.Fatal("Telegram bot token is missing in configuration file")
	}

	if dbUrl == "" {
		log.Fatal("Database url is missing in configuration file")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	db, err := database.PostgresPool(dbUrl)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}
	defer db.Close()

	log.Println("Database connected")

	botClient, err := bot.New(token)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepository(db)

	services := &app.Services{
		Admin:  *service.NewAdminService(repo, botClient),
		Config: *service.NewConfigService(repo),
		Chat:   *service.NewChatService(repo, botClient),
	}

	mw := middleware.New(services.Chat)

	app := app.New(db, services)

	allCommands := []*command.Command{
		factcheck.New(&app.GeminiService),
		summary.New(repo, &app.GeminiService),
		ask.New(&app.GeminiService),
		look.New(&app.GeminiService),
		duplicate.New(),
		status.New(),
	}

	cm := command.New()
	cm.RegisterAll(allCommands)

	duplicateFeature := duplicate.NewDuplicate(&services.Chat)
	duplicateMenuHandler := duplicate.NewMenuHandler(duplicateFeature)
	menuManager := menu.NewMenuManager(duplicateMenuHandler, cm, app)

	r := router.NewRouter()

	r.Register("instagram", embedHandler.Instagram)
	r.Register("tiktok", embedHandler.TikTok)

	botClient.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact,
		func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			menuManager.HandleStart(ctx, bot, update)
		})

	botClient.RegisterHandlerMatchFunc(
		func(update *models.Update) bool {
			return update != nil &&
				update.Message != nil &&
				update.Message.Chat.Type == models.ChatTypePrivate &&
				update.Message.Text != "/start"
		},
		func(ctx context.Context, b *bot.Bot, update *models.Update) {
			menuManager.HandleMessage(ctx, b, update)
		},
	)

	botClient.RegisterHandlerMatchFunc(
		func(update *models.Update) bool {
			return update != nil &&
				update.Message != nil &&
				update.Message.Chat.Type != models.ChatTypePrivate
		},
		dispatcher.MainHandler(app, r, mw, cm))

	botClient.RegisterHandlerMatchFunc(
		func(update *models.Update) bool {
			return update != nil && update.CallbackQuery != nil
		},
		func(ctx context.Context, b *bot.Bot, update *models.Update) {
			menuManager.HandleCallback(ctx, b, update)
		})

	botClient.RegisterHandlerMatchFunc(func(update *models.Update) bool {
		logger.DebugJson(update)
		if update == nil || update.MyChatMember == nil {
			return false
		}

		myChatMember := update.MyChatMember
		botID := botClient.ID()

		if myChatMember.NewChatMember.Member != nil &&
			myChatMember.NewChatMember.Member.User != nil &&
			myChatMember.NewChatMember.Member.User.ID == botID {
			return true
		}

		if myChatMember.NewChatMember.Left != nil &&
			myChatMember.NewChatMember.Left.User != nil &&
			myChatMember.NewChatMember.Left.User.ID == botID {
			return true
		}

		return false
	},
		func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			botID := bot.ID()
			myChatMember := update.MyChatMember

			newMember := myChatMember.NewChatMember

			if newMember.Member != nil &&
				newMember.Member.User != nil &&
				newMember.Member.User.IsBot &&
				newMember.Member.User.ID == botID {

				chatHandler.HandleJoinChat(ctx, bot, update, app)
			}

			if newMember.Left != nil &&
				newMember.Left.User != nil &&
				newMember.Left.User.IsBot &&
				newMember.Left.User.ID == botID {

				chatHandler.HandleLeaveChat(ctx, bot, update, app.DB)
			}
		})

	botClient.RegisterHandlerMatchFunc(func(update *models.Update) bool {
		return update != nil && update.EditedMessage != nil
	}, func(ctx context.Context, bot *bot.Bot, update *models.Update) {
		messageHandler.UpdateMessage(ctx, bot, update, db)
	})

	log.Println("Bot started")

	botClient.Start(ctx)
}
