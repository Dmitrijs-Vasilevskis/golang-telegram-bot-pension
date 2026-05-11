package duplicate

import (
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/command"
)

func New() *command.Command {
	return &command.Command{
		Name:        "Duplicate DM",
		Key:         "duplicate_dm",
		Description: "Sends a copy of messages you send to the bot in private (DM) to the selected group chat.",
		Icon:        "💾",
		IsFeature:   true,
		System:      false,
	}
}
