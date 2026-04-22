package middleware

import (
	"context"
)

func (m *Middleware) FeatureMiddleware(ctx context.Context, chatID int64, userID, feature string) (bool, error) {
	var enabled bool
	var err error

	// if feature == "" {
	// 	return false, fmt.Errorf("")
	// }

	switch feature {
	case "sumary":
		enabled, err = m.SummaryMiddleware(ctx, chatID)
	case "duplicate":
		enabled, err = m.DuplicateMiddleware(ctx, chatID)
	}

	return enabled, err
}

func (m *Middleware) DuplicateMiddleware(ctx context.Context, chatID int64) (bool, error) {
	enabled, err := m.config.IsDuplicateDMEnabled(ctx, chatID)
	if err != nil {
		return false, err
	}

	return enabled, nil
}

func (m *Middleware) SummaryMiddleware(ctx context.Context, chatID int64) (bool, error) {
	enabled, err := m.config.IsSummaryEnabled(ctx, chatID)
	if err != nil {
		return false, err
	}

	return enabled, nil
}
