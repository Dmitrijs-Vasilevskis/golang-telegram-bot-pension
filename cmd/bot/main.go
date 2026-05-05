package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/app"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/database"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/dispatcher"
	chatHandler "github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/handlers/chat"
	embedHandler "github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/handlers/embed"
	messageHandler "github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/handlers/messages"
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
	menuManager := menu.NewMenuManager(app)

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
		dispatcher.MainHandler(app, r, mw))

	botClient.RegisterHandlerMatchFunc(
		func(update *models.Update) bool {
			return update != nil && update.CallbackQuery != nil
		},
		func(ctx context.Context, b *bot.Bot, update *models.Update) {
			menuManager.HandleCallback(ctx, b, update)
		})

	botClient.RegisterHandlerMatchFunc(func(update *models.Update) bool {
		if update == nil || update.MyChatMember == nil {
			return false
		}

		myChatMember := update.MyChatMember

		if myChatMember != nil &&
			myChatMember.NewChatMember.Member != nil &&
			myChatMember.NewChatMember.Member.User.ID == botClient.ID() {
			return true
		}

		if myChatMember != nil &&
			myChatMember.OldChatMember.Member != nil &&
			myChatMember.OldChatMember.Left.User.ID == botClient.ID() {
			return true
		}

		return false
	},
		func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			botID := bot.ID()
			myChatMember := update.MyChatMember

			newMember := myChatMember.NewChatMember

			if newMember.Type == models.ChatMemberTypeMember && newMember.Member.User.IsBot &&
				(newMember.Member.User.ID == botID) {
				chatHandler.HandleJoinChat(ctx, bot, update, app)
			}

			if newMember.Type == models.ChatMemberTypeLeft && newMember.Left.User.IsBot &&
				(newMember.Left.User.ID == botID) {
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
