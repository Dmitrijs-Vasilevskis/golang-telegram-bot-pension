package middleware

import (
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/service"
)

type Middleware struct {
	config service.ChatService
}

func New(config service.ChatService) *Middleware {
	return &Middleware{
		config: config,
	}
}
