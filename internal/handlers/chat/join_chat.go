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

	// Проверяем, что это обновление о боте
	newMember := myChatMember.NewChatMember
	if newMember.Member == nil || newMember.Member.User == nil {
		return
	}

	// Проверяем, что бот был добавлен
	if !newMember.Member.User.IsBot || newMember.Member.User.ID != b.ID() {
		return
	}

	// Проверяем статус "member" (бот добавлен)
	if newMember.Member.Status != "member" {
		return
	}

	chat := myChatMember.Chat
	userWhoAdded := myChatMember.From

	log.Printf("Bot added to chat %d (%s) by user %d", chat.ID, chat.Title, userWhoAdded.ID)

	// Регистрируем чат в системе
	chatService := app.Services.Chat
	if err := chatService.RegisterChat(ctx, chat.ID, chat.Title, string(chat.Type), &userWhoAdded); err != nil {
		log.Printf("Failed to register chat %d: %v", chat.ID, err)

		// Отправляем сообщение об ошибке
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat.ID,
			Text:   "❌ Failed to initialize bot settings. Please remove and add the bot again.",
		})
		return
	}

	// Отправляем приветственное сообщение с настройками по умолчанию
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
