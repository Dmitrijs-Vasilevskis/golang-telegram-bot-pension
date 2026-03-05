package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/handlers"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/helpers"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/logger"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/router"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Cannot load the env file")
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")

	if token == "" {
		log.Fatal("Telegram bot token is missing in configuration file")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	botClient, err := bot.New(token)
	if err != nil {
		log.Fatal(err)
	}

	r := router.NewRouter()

	r.Register("reel", handlers.Reels)
	r.Register("tiktok", handlers.TikTok)

	botClient.RegisterHandler(bot.HandlerTypeMessageText, "/status", bot.MatchTypeExact,
		func(ctx context.Context, botClient *bot.Bot, update *models.Update) {
			botClient.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Bot is Active 🚀",
			})
		})

	botClient.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeContains,
		func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			text := update.Message.Text

			logger.LogUpdate(update)

			if helpers.IsToxic(text) {
				handlers.Clown(ctx, bot, update)
			}

			r.Handle(ctx, bot, update)
		})

	log.Println("Bot started")

	botClient.Start(ctx)
}
