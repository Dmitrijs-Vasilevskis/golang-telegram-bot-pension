package handlers

import (
	"context"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Duplicator(ctx context.Context, b *bot.Bot, update *models.Update) {
	msg := update.Message
	targetChatID := -1001304243920
	var err error

	switch {
	case msg.Text != "":
		text := msg.Text
		if strings.Contains(text, ")") {
			text = text + " 🤡"
		}
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: targetChatID,
			Text:   text,
		})
	case len(msg.Photo) > 0:
		photo := msg.Photo[len(msg.Photo)-1]
		_, err = b.SendPhoto(ctx, &bot.SendPhotoParams{
			ChatID:  targetChatID,
			Photo:   &models.InputFileString{Data: photo.FileID},
			Caption: msg.Caption,
		})
	case msg.Video != nil:
		_, err = b.SendVideo(ctx, &bot.SendVideoParams{
			ChatID:  targetChatID,
			Video:   &models.InputFileString{Data: msg.Video.FileID},
			Caption: msg.Caption,
		})
	case msg.Voice != nil:
		_, err = b.SendVoice(ctx, &bot.SendVoiceParams{
			ChatID:  targetChatID,
			Voice:   &models.InputFileString{Data: msg.Voice.FileID},
			Caption: msg.Caption,
		})
	case msg.Sticker != nil:
		_, err = b.SendSticker(ctx, &bot.SendStickerParams{
			ChatID:  targetChatID,
			Sticker: &models.InputFileString{Data: msg.Sticker.FileID},
		})
	case msg.VideoNote != nil:
		_, err = b.SendVideoNote(ctx, &bot.SendVideoNoteParams{
			ChatID:    targetChatID,
			VideoNote: &models.InputFileString{Data: msg.VideoNote.FileID},
		})
	case msg.Audio != nil:
		_, err = b.SendAudio(ctx, &bot.SendAudioParams{
			ChatID:  targetChatID,
			Audio:   &models.InputFileString{Data: msg.Audio.FileID},
			Caption: msg.Caption,
		})
	}

	if err != nil {
		log.Printf("Failed to forward DM: %v", err)
	} else {
		log.Printf("Forwarded DM from user %d to chat %d", msg.From.ID, targetChatID)
	}
}
