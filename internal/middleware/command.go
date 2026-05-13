package middleware

import (
	"context"
)

func (m *Middleware) CommandMiddleware(ctx context.Context, chatID int64, command string) (bool, error) {
	enabled, err := m.config.IsCommandEnabled(ctx, chatID, command)
	if err != nil {
		return false, err
	}

	return enabled, nil
}
