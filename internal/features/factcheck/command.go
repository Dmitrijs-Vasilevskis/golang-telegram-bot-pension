package factcheck

import (
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/command"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/service"
)

func New(gemini *service.GeminiService) *command.Command {
	return &command.Command{
		Name:        "Facksheck",
		Key:         "factcheck",
		Description: "Verify whether a statement is true or misleading.",
		IsFeature:   false,
		System:      false,
		Handler:     Handler(gemini),
	}
}
