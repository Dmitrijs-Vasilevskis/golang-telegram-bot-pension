package status

import "github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/command"

func New() *command.Command {
	return &command.Command{
		Name:        "Status",
		Key:         "status",
		Description: "Check current bot status",
		IsFeature:   false,
		System:      true,
		Handler:     StatusHandler(),
	}
}
