package text

import (
	"MemoryPicBot/bot/repository"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/sirupsen/logrus"
)

type PicDeleter interface {
	DelPic(userID int64, tag string) error
}

// DelPic handles text messages in the tag deletion flow, removes the picture with the specified tag from both database and cloud storage
func (h *HandlerText) DelPic(message *tgbotapi.Message) {
	tag := message.Text

	picKey, err := h.repo.DelPic(message.From.ID, tag)
	if err != nil {
		if errors.Is(err, repository.PicNotFound) {
			logrus.Warnf("pictures not found in your collection for tag: %s", tag)

			h.bot.SendMessage(message.From.ID, fmt.Sprintf("%v:%v", repository.PicNotFound, tag))

			return
		}
		h.bot.HandleError(message.From.ID, h.bot.Cfg.Msg.ErrDelPic, "delete pic failed from repo: %v", err)

		return
	}

	if err := h.cloud.DeletePic(picKey); err != nil {
		h.bot.HandleError(message.From.ID, h.bot.Cfg.Msg.ErrDelPic, "delete pic failed from s3Cloud: %v", err)

		return
	}

	h.bot.OKResponse(message, h.bot.Cfg.Msg.PicDeleted)
}
