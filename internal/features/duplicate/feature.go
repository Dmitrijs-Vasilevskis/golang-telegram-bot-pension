package duplicate

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
)

type FeatureStore interface {
	IsDuplicateDMEnabled(ctx context.Context, chatID int64) (bool, error)
}

type Feature struct {
	store FeatureStore
}

func NewDuplicate(store FeatureStore) *Feature {
	return &Feature{store: store}
}

func (f *Feature) CopyMessage(
	ctx context.Context,
	b *bot.Bot,
	toChatID int64,
	fromChatID int64,
	messageID int,
	dmChannelID int64,
) error {
	enabled, err := f.store.IsDuplicateDMEnabled(ctx, toChatID)

	if err != nil {
		return fmt.Errorf("failed to check duplicate feature status: %w", err)
	}

	if !enabled {
		return fmt.Errorf("duplicate feature is not enabled for chat %d", fromChatID)
	}

	_, err = b.CopyMessage(ctx, &bot.CopyMessageParams{
		ChatID:     toChatID,
		FromChatID: fromChatID,
		MessageID:  messageID,
	})

	return err
}
