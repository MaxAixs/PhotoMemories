package text

import (
	"MemoryPicBot/bot"
	"MemoryPicBot/bot/handlers/state"
	"MemoryPicBot/bot/repository"
	"MemoryPicBot/pkg/cache"
	"MemoryPicBot/pkg/s3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HandlerText struct {
	bot       *bot.Bot
	repo      *repository.PicRepository
	cache     *cache.Client
	cloud     *s3.Client
	userState *state.Manager
}

func NewTextHandler(bot *bot.Bot, repo *repository.PicRepository, redis *cache.Client, s3 *s3.Client, userState *state.Manager) *HandlerText {
	return &HandlerText{
		bot:       bot,
		repo:      repo,
		cache:     redis,
		cloud:     s3,
		userState: userState,
	}
}

// DoText handles user messages based on their current state.
func (h *HandlerText) DoText(message *tgbotapi.Message) {
	userState := h.userState.GetUserState(message.From.ID)

	switch userState {
	case state.AwaitSaveTag:
		h.SavePic(message)

	case state.AwaitDelTag:
		h.DelPic(message)

	case state.AwaitGetTag:
		h.GetPic(message)

	case state.Default:
		h.Default(message)

	}
}

// Default sends a default response when the input doesn't match any known commands.
func (h *HandlerText) Default(message *tgbotapi.Message) {
	h.bot.SendMessage(message.From.ID, h.bot.Cfg.Msg.Default)
}
