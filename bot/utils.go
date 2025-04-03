package bot

import (
	"MemoryPicBot/bot/handlers/state"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func (b *Bot) SendMessage(chatId int64, message string) {
	msg := tgbotapi.NewMessage(chatId, message)

	b.TgAPI.Send(msg)
}

func (b *Bot) OKResponse(message *tgbotapi.Message, response string) {
	b.SendMessage(message.From.ID, response)
	b.Manager.SetState(message.From.ID, state.Default)
	b.SendInLineCmd(message.From.ID)
}

func (b *Bot) CallbackResponse(callback *tgbotapi.CallbackQuery, response string, state string) {
	b.SendMessage(callback.Message.Chat.ID, response)
	b.Manager.SetState(callback.From.ID, state)
}

func (b *Bot) HandleError(UserID int64, response string, logMsg string, err error) {
	logrus.Errorf(logMsg, err)
	b.SendMessage(UserID, response)
}
