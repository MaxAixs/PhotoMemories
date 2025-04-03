package text

import (
	"MemoryPicBot/bot/repository"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type PicGetter interface {
	GetPic(userID int64, tag string) (string, error)
}

// GetPic retrieves and sends a picture to the user based on the provided tag
func (h *HandlerText) GetPic(message *tgbotapi.Message) {
	tag := message.Text
	picKey, err := h.repo.GetPic(message.From.ID, tag)
	if err != nil {
		if errors.Is(err, repository.PicNotFound) {
			logrus.Warnf("Pictures not found in your collection for tag: %s", tag)

			h.bot.SendMessage(message.From.ID, fmt.Sprintf("%v:%v", repository.PicNotFound, tag))

			return
		}
		h.bot.HandleError(message.From.ID, h.bot.Cfg.Msg.ErrGetPic, "Get pictures failed %v", err)

		return
	}

	uploadPic, err := h.cloud.DownLoadPic(picKey)
	if err != nil {
		h.bot.HandleError(message.From.ID, h.bot.Cfg.Msg.ErrGetPic, "Download pictures from s3Cloud failed %v", err)

		return
	}

	h.sendPic(message.From.ID, uploadPic, tag)
	h.bot.OKResponse(message, tag)
}

// sendPic sends a photo to the specified chat with the given tag as the photo name
func (h *HandlerText) sendPic(chatID int64, pic []byte, tag string) {
	pictures := tgbotapi.NewPhoto(chatID, tgbotapi.FileBytes{
		Name:  tag,
		Bytes: pic,
	})

	h.bot.TgAPI.Send(pictures)
}
