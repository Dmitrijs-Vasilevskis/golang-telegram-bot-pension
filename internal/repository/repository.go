package repository

import "github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/database"

type Repository struct {
	db database.DBTX
}

func NewRepository(db database.DBTX) *Repository {
	return &Repository{db}
}
