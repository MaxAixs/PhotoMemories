package bot

import (
	"MemoryPicBot/bot/config"
	"MemoryPicBot/bot/handlers/state"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	TgAPI   *tgbotapi.BotAPI
	Cfg     *config.Config
	Buttons *tgbotapi.InlineKeyboardMarkup
	Manager *state.Manager
}

func NewBot(tgBotAPI *tgbotapi.BotAPI, cfg *config.Config, userState *state.Manager) *Bot {
	return &Bot{
		TgAPI:   tgBotAPI,
		Cfg:     cfg,
		Buttons: createInlineKeyboard(cfg.Buttons),
		Manager: userState,
	}
}

type UpdateProcessor interface {
	ProcessUpdate(update tgbotapi.Update)
}

func (b *Bot) Start(processor UpdateProcessor) {
	logrus.Printf("Authorized on account UserName: %s, UserID: %v", b.TgAPI.Self.UserName, b.TgAPI.Self.ID)

	updates := b.GetUpdates()

	for update := range updates {
		processor.ProcessUpdate(update)
	}
}

func (b *Bot) GetUpdates() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.TgAPI.GetUpdatesChan(u)

	return updates
}
