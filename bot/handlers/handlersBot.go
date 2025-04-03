package handlers

import (
	tgbot "MemoryPicBot/bot"
	"MemoryPicBot/bot/handlers/buttons"
	"MemoryPicBot/bot/handlers/cmd"
	"MemoryPicBot/bot/handlers/pictures"
	"MemoryPicBot/bot/handlers/state"
	"MemoryPicBot/bot/handlers/text"
	"MemoryPicBot/bot/repository"
	"MemoryPicBot/pkg/cache"
	"MemoryPicBot/pkg/s3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotHandlers struct {
	Btn *buttons.HandlerBtn
	Cmd *cmd.HandlerCmd
	Pic *pictures.HandlerPic
	Txt *text.HandlerText
}

func InitHandlers(bot *tgbot.Bot, redis *cache.Client, repo *repository.PicRepository, s3Cloud *s3.Client, userState *state.Manager) *BotHandlers {

	return &BotHandlers{
		Btn: buttons.NewBotHandlerBtn(bot, repo, userState),
		Cmd: cmd.NewHandleCmd(bot),
		Pic: pictures.NewPicHandler(bot, redis, s3Cloud, userState),
		Txt: text.NewTextHandler(bot, repo, redis, s3Cloud, userState),
	}
}

func (b *BotHandlers) ProcessUpdate(update tgbotapi.Update) {
	switch {
	case update.CallbackQuery != nil:
		b.Btn.DoCallbackQuery(update.CallbackQuery)

	case update.Message == nil:
		return

	case update.Message.IsCommand():
		b.Cmd.DoCmd(update.Message)

	case update.Message.Photo != nil && len(update.Message.Photo) > 0:
		b.Pic.DoPictures(update.Message)

	case update.Message.Text != "":
		b.Txt.DoText(update.Message)

	}
}
