// internal/handlers/chat.go
package handlers

import (
	"context"
	"log"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/app"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func HandleJoinChat(ctx context.Context, b *bot.Bot, update *models.Update, app *app.App) {
	myChatMember := update.MyChatMember
	if myChatMember == nil {
		return
	}

	newMember := myChatMember.NewChatMember
	if newMember.Member == nil || newMember.Member.User == nil {
		return
	}

	if !newMember.Member.User.IsBot || newMember.Member.User.ID != b.ID() {
		return
	}

	if newMember.Member.Status != "member" {
		return
	}

	chat := myChatMember.Chat
	userWhoAdded := myChatMember.From

	log.Printf("Bot added to chat %d (%s) by user %d", chat.ID, chat.Title, userWhoAdded.ID)

	chatService := app.Services.Chat
	if err := chatService.RegisterChat(ctx, chat.ID, chat.Title, string(chat.Type), &userWhoAdded); err != nil {
		log.Printf("Failed to register chat %d: %v", chat.ID, err)

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat.ID,
			Text:   "❌ Failed to initialize bot settings. Please remove and add the bot again.",
		})
		return
	}

	welcomeMsg := "🤖 *Bot activated!*\n\n" +
		"Default settings:\n" +
		"• 📝 Summary: ❌ Disabled\n" +
		"• 🔄 Duplicate DM: ❌ Disabled\n" +
		"• 🎮 Commands: ✅ All enabled\n\n" +
		"Use `/settings` in private chat with me to configure features."

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chat.ID,
		Text:      welcomeMsg,
		ParseMode: models.ParseModeMarkdown,
	})
}
