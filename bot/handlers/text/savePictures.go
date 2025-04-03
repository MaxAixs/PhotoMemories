package text

import (
	"MemoryPicBot/bot/repository"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PicSaver interface {
	SavePic(userID int64, picKey string, tag string) error
}

// SavePic saves the uploaded photo of the user in the database with the specified tag.
func (h *HandlerText) SavePic(message *tgbotapi.Message) {
	tag := message.Text

	picKey, err := h.cache.GetPicture(message.From.ID)
	if err != nil {
		h.bot.HandleError(message.From.ID, h.bot.Cfg.Msg.ErrUploadPic, "cant get pictures from Redis: %v", err)
		return
	}

	if err := h.repo.SavePic(message.From.ID, picKey, tag); err != nil {
		if errors.Is(err, repository.PicExists) {
			h.bot.SendMessage(message.From.ID, repository.PicExists.Error())
			return
		}
		h.bot.HandleError(message.From.ID, h.bot.Cfg.Msg.ErrUploadTag, "save pic failed: %v", err)
		return
	}

	h.bot.OKResponse(message, h.bot.Cfg.Msg.TagSaved)
}
