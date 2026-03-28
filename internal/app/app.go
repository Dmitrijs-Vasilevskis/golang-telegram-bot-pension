package app

import (
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/repository"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Services struct {
	Admin  service.AdminService
	Config service.ConfigService
	Chat   service.ChatService
}

type App struct {
	DB *pgxpool.Pool

	Repository    *repository.Repository
	GeminiService service.GeminiService
	Services      *Services
}

func New(db *pgxpool.Pool, services *Services) *App {
	return &App{
		DB: db,

		Repository:    repository.NewRepository(db),
		GeminiService: *service.NewGeminiService(),
		Services:      services,
	}
}
