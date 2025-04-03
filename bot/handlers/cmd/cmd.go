package cmd

import (
	"MemoryPicBot/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type HandlerCmd struct {
	bot *bot.Bot
}

func NewHandleCmd(bot *bot.Bot) *HandlerCmd {
	return &HandlerCmd{
		bot: bot,
	}
}

// DoCmd handles command messages (those starting with /) and routes them to the appropriate handler function
func (h *HandlerCmd) DoCmd(updateMsg *tgbotapi.Message) {
	switch updateMsg.Command() {
	case h.bot.Cfg.Cmd.Start:
		logrus.Infof("Received /start command from user: %v", updateMsg.Chat.ID)

		h.cmdStart(updateMsg.Chat.ID)

	case h.bot.Cfg.Cmd.Help:
		logrus.Infof("Received /help command from user: %v", updateMsg.Chat.ID)

		h.cmdHelp(updateMsg.Chat.ID)
	}

}

// cmdStart processes the /start command, sending the welcome message and displaying inline commands to the user
func (h *HandlerCmd) cmdStart(chatID int64) {
	h.bot.SendMessage(chatID, h.bot.Cfg.Msg.Start)
	h.bot.SendInLineCmd(chatID)
}

// cmdHelp processes the /help command, sending help information to the user
func (h *HandlerCmd) cmdHelp(chatID int64) {
	h.bot.SendMessage(chatID, h.bot.Cfg.Msg.Help)
}
