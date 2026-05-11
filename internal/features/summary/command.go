package summary

import (
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/command"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/repository"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/service"
)

func New(repo *repository.Repository, gemini *service.GeminiService) *command.Command {
	return &command.Command{
		Name: "Summary",
		Key:  "summary",
		Description: "Quickly generate a summary of recent chat activity.\n\n" +
			"Use `/summary <count>` (e.g. `/summary 200`) to choose how many messages to include.\n\n" +
			"• Maximum: 400 messages\n" +
			"• Messages are anonymized and securely processed",
		Icon:      "💾",
		IsFeature: true,
		System:    false,
		Handler:   Handler(repo, gemini),
	}
}
