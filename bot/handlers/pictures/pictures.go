package pictures

import (
	"MemoryPicBot/bot"
	"MemoryPicBot/bot/handlers/state"
	"MemoryPicBot/pkg/cache"
	"MemoryPicBot/pkg/s3"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"io"
	"net/http"
)

type HandlerPic struct {
	bot       *bot.Bot
	cache     *cache.Client
	cloud     *s3.Client
	userState *state.Manager
}

func NewPicHandler(bot *bot.Bot, redis *cache.Client, s3Client *s3.Client, userState *state.Manager) *HandlerPic {
	return &HandlerPic{
		bot:       bot,
		cache:     redis,
		cloud:     s3Client,
		userState: userState,
	}
}

// DoPictures processes photo uploads from users in the AwaitPic state, downloads the photo, uploads it to cloud storage, and updates user state
func (h *HandlerPic) DoPictures(message *tgbotapi.Message) {
	if !h.userState.IsUserInState(message.From.ID, state.AwaitPic) {
		return
	}

	if len(message.Photo) == 0 {
		return
	}

	picID := extractPhotoID(message)
	file, err := h.bot.TgAPI.GetFile(tgbotapi.FileConfig{FileID: picID})
	if err != nil {
		h.bot.HandleError(message.From.ID, h.bot.Cfg.Msg.ErrUploadPic, "Error getting photo file from telegram: %v", err)

		return
	}

	picData, err := h.downloadPic(file)
	if err != nil {
		h.bot.HandleError(message.From.ID, h.bot.Cfg.Msg.ErrUploadPic, "Error downloading photo file: %v", err)

		return
	}

	picKey := generatePicKey()
	err = h.cloud.UploadPic(picKey, picData)
	if err != nil {
		h.bot.HandleError(message.From.ID, h.bot.Cfg.Msg.ErrUploadPic, "Error uploading photo to S3: %v", err)

		return
	}

	if err := h.cache.SavePicture(message.From.ID, picKey); err != nil {
		h.bot.HandleError(message.From.ID, h.bot.Cfg.Msg.ErrUploadPic, "Error saving pictures: %v", err)
		return
	}

	h.bot.SendMessage(message.From.ID, h.bot.Cfg.Msg.PicSaved)
	h.userState.SetState(message.From.ID, state.AwaitSaveTag)
}

// extractPhotoID retrieves the FileID of the highest resolution photo from the message
func extractPhotoID(message *tgbotapi.Message) string {
	return message.Photo[len(message.Photo)-1].FileID
}

// generatePicKey creates a unique UUID identifier for storing the photo
func generatePicKey() string {
	uuidStr := uuid.New().String()

	return uuidStr
}

// downloadPic retrieves a photo file from Telegram servers using the provided file information
func (h *HandlerPic) downloadPic(file tgbotapi.File) ([]byte, error) {
	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", h.bot.Cfg.TelegramBotToken, file.FilePath)

	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("error download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error download file: status code %d", resp.StatusCode)
	}

	photoData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	return photoData, nil
}
