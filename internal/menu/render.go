package menu

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (mm *MenuManager) renderCommandEdit(ctx context.Context, b *bot.Bot, state *MenuState) {
	cmdKey := state.Data["command"].(string)

	cmd, _ := mm.commandManager.GetByKey(cmdKey)

	enabled, _ := mm.chatService.IsCommandEnabled(ctx, state.ChatID, cmd.Key)

	status := "❌ Disabled"
	if enabled {
		status = "✅ Enabled"
	}

	keyboard := mm.actionKeyboard(cmd.Key)

	text := fmt.Sprintf(
		"%s *%s*\n\n*Description:*\n%s\n\n*Status:* %s",
		cmd.Icon,
		cmd.Name,
		cmd.Description,
		status,
	)

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      state.DMChannelID,
		MessageID:   state.MessageId,
		ParseMode:   models.ParseModeMarkdownV1,
		Text:        text,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		fmt.Printf("Failed to render a commands edit keyboard: %v\n", err)
	}
}

func (mm *MenuManager) renderCommands(ctx context.Context, b *bot.Bot, state *MenuState) {
	keyboard := mm.buildCommandsKeyboard(ctx, state)

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      state.DMChannelID,
		MessageID:   state.MessageId,
		ParseMode:   models.ParseModeMarkdown,
		Text:        "Command Management",
		ReplyMarkup: keyboard,
	})
	if err != nil {
		fmt.Printf("Failed to render a commands keyboard: %v\n", err)
	}
}

func (mm *MenuManager) renderFeatures(ctx context.Context, b *bot.Bot, state *MenuState) {
	keyboard := mm.buildFeaturesKeyboard()

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      state.DMChannelID,
		MessageID:   state.MessageId,
		ParseMode:   models.ParseModeMarkdown,
		Text:        "Select a feature to configure",
		ReplyMarkup: keyboard,
	})
	if err != nil {
		fmt.Printf("Failed to render a features keyboard: %v\n", err)
	}
}

func (mm *MenuManager) renderFeatureEdit(ctx context.Context, b *bot.Bot, state *MenuState) {
	featureKey := state.Data["feature"].(string)

	feature, _ := mm.commandManager.GetByKey(featureKey)

	enabled, _ := mm.chatService.IsFeatureEnabled(ctx, state.ChatID, feature.Key)

	status := "❌ Disabled"
	if enabled {
		status = "✅ Enabled"
	}

	keyboard := mm.buildFeatureEditKeyboard(feature.Key)

	text := fmt.Sprintf(
		"%s *%s*\n\n*Description:*\n%s\n\n*Status:* %s",
		feature.Icon,
		feature.Name,
		feature.Description,
		status,
	)

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      state.DMChannelID,
		MessageID:   state.MessageId,
		ParseMode:   models.ParseModeMarkdownV1,
		Text:        text,
		ReplyMarkup: keyboard,
	})

	if err != nil {
		fmt.Printf("Failed to render a feature edit keyboard: %v\n", err)
	}
}

func (mm *MenuManager) renderMain(ctx context.Context, b *bot.Bot, state *MenuState) {
	keyboard := mm.mainKeyboard()

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      state.DMChannelID,
		MessageID:   state.MessageId,
		ParseMode:   models.ParseModeMarkdown,
		Text:        "*Main Menu*",
		ReplyMarkup: keyboard,
	})
	if err != nil {
		fmt.Printf("Failed to render a main keyboard: %v\n", err)
	}
}

func (mm *MenuManager) renderChats(ctx context.Context, b *bot.Bot, state *MenuState) {
	keyboard := mm.buildChatsKeyboard(ctx, *state)

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      state.DMChannelID,
		MessageID:   state.MessageId,
		ParseMode:   models.ParseModeMarkdown,
		Text:        "Select a chat:",
		ReplyMarkup: keyboard,
	})
	if err != nil {
		fmt.Printf("Failed to render a chats keyboard: %v\n", err)
	}
}

func (mm *MenuManager) renderChatActions(ctx context.Context, b *bot.Bot, state *MenuState) {
	keyboard := mm.buildChatActionsBeyboard()

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      state.DMChannelID,
		MessageID:   state.MessageId,
		ParseMode:   models.ParseModeMarkdown,
		Text:        "Select an action:",
		ReplyMarkup: keyboard,
	})
	if err != nil {
		fmt.Printf("Failed to render a setting keyboard: %v\n", err)
	}
}

func (mm *MenuManager) renderSettings(ctx context.Context, b *bot.Bot, state *MenuState) {
	keyboard := mm.buildSettingKeyboard()

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      state.DMChannelID,
		MessageID:   state.MessageId,
		ParseMode:   models.ParseModeMarkdown,
		Text:        "Select a feature to configure:",
		ReplyMarkup: keyboard,
	})
	if err != nil {
		fmt.Printf("Failed to render a setting keyboard: %v\n", err)
	}
}

func (mm *MenuManager) renderDuplicateAction(ctx context.Context, b *bot.Bot, state *MenuState) {
	keyboard := mm.buildDuplicateActionKeyboard()

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      state.DMChannelID,
		MessageID:   state.MessageId,
		ParseMode:   models.ParseModeMarkdownV1,
		Text:        "✉️ Send a message to duplicate.\n\nType /cancel to abort.",
		ReplyMarkup: keyboard,
	})

	if err != nil {
		fmt.Printf("Failed to render a keyboard: %v\n", err)
	}
}
