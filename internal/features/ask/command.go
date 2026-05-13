package ask

import (
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/command"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/service"
)

func New(gemini *service.GeminiService) *command.Command {
	return &command.Command{
		Name:        "Ask",
		Key:         "ask",
		Description: "Ask the bot any question and get an AI-powered answer.",
		IsFeature:   false,
		System:      false,
		Handler:     Handler(gemini),
	}
}
