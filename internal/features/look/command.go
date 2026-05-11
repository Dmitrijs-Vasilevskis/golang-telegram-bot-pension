package look

import (
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/command"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/service"
)

func New(gemini *service.GeminiService) *command.Command {
	return &command.Command{
		Name:        "Look",
		Key:         "look",
		Description: "Analyze messages or media content in the chat.",
		IsFeature:   false,
		System:      false,
		Handler:     Handler(gemini),
	}
}
