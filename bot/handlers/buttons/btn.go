package buttons

import (
	"MemoryPicBot/bot"
	"MemoryPicBot/bot/handlers/state"
	"MemoryPicBot/bot/repository"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"strings"
)

type HandlerBtn struct {
	bot       *bot.Bot
	repo      *repository.PicRepository
	userState *state.Manager
}

func NewBotHandlerBtn(bot *bot.Bot, repo *repository.PicRepository, userState *state.Manager) *HandlerBtn {
	return &HandlerBtn{bot: bot, userState: userState, repo: repo}
}

type TagsGetter interface {
	GetAllTags(userID int64) ([]string, error)
}

// DoCallbackQuery handles incoming callback queries from Telegram buttons
func (b *HandlerBtn) DoCallbackQuery(callback *tgbotapi.CallbackQuery) {
	switch callback.Data {
	case bot.BtnAddPhoto:
		b.btnAddPic(callback)

	case bot.BtnDelete:
		b.btnDelPic(callback)

	case bot.BtnGetPhoto:
		b.btnGetPic(callback)

	case bot.BtnMyTags:
		b.btnGetTags(callback)
	}
}

// btnAddPic handles the "Add Photo" button click, sends instructions to the user and sets state to await photo upload
func (b *HandlerBtn) btnAddPic(callback *tgbotapi.CallbackQuery) {
	b.bot.CallbackResponse(callback, b.bot.Cfg.Msg.AddPic, state.AwaitPic)
}

// btnDelPic handles the "Delete Photo" button click, sends instructions to the user and sets state to await tag for deletion
func (b *HandlerBtn) btnDelPic(callback *tgbotapi.CallbackQuery) {
	b.bot.CallbackResponse(callback, b.bot.Cfg.Msg.DelPic, state.AwaitDelTag)
}

// btnGetPic handles the "Get Photo" button click, sends instructions to the user and sets state to await tag for retrieval
func (b *HandlerBtn) btnGetPic(callback *tgbotapi.CallbackQuery) {
	b.bot.CallbackResponse(callback, b.bot.Cfg.Msg.GetPic, state.AwaitGetTag)
}

// btnGetTags handles the "My Tags" button click and calls the getTags function to display user's saved tags
func (b *HandlerBtn) btnGetTags(callback *tgbotapi.CallbackQuery) {
	b.getTags(callback)
}

// getTags retrieves and displays all tags associated with the current user, handling various error cases
func (b *HandlerBtn) getTags(callback *tgbotapi.CallbackQuery) {
	tags, err := b.repo.GetAllTags(callback.From.ID)
	if errors.Is(err, repository.UserNotFound) {
		logrus.Infof("User not found %v", callback.Message.Chat.ID)

		b.bot.SendMessage(callback.Message.Chat.ID, repository.UserNotFound.Error())

		return
	}
	if err != nil {
		b.bot.HandleError(callback.From.ID, b.bot.Cfg.Msg.ErrGetTags, "Get tags failed: %v", err)

		return
	}

	b.bot.SendMessage(callback.Message.Chat.ID, fmt.Sprintf("Your tags: %v", strings.Join(tags, ",")))

}
