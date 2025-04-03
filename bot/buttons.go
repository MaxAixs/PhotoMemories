package bot

import (
	"MemoryPicBot/bot/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	BtnAddPhoto = "btn_add_photo"
	BtnDelete   = "btn_delete_photo"
	BtnGetPhoto = "btn_get_photo"
	BtnMyTags   = "btn_my_tags"
)

func createInlineKeyboard(cfg config.Buttons) *tgbotapi.InlineKeyboardMarkup {
	btn := initButtons(cfg)

	return &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: btn,
	}
}

func initButtons(cfg config.Buttons) [][]tgbotapi.InlineKeyboardButton {
	buttons := []struct {
		text string
		data string
	}{
		{cfg.AddPic, BtnAddPhoto}, {cfg.DelPic, BtnDelete},
		{cfg.GetPic, BtnGetPhoto}, {cfg.MyTags, BtnMyTags},
	}

	var inlineButtons [][]tgbotapi.InlineKeyboardButton
	for _, btn := range buttons {
		inlineButtons = append(inlineButtons, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(btn.text, btn.data),
		})
	}

	return inlineButtons
}

func (b *Bot) SendInLineCmd(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Available actions:")

	msg.ReplyMarkup = b.Buttons

	b.TgAPI.Send(msg)
}
