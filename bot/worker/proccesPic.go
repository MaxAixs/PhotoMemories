package worker

import (
	"MemoryPicBot/bot"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"sync"
)

func (w *Worker) ProcessPicLists(picList []bot.UserPictures) []error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(picList))

	for _, item := range picList {
		wg.Add(1)

		go func(item bot.UserPictures) {
			defer wg.Done()
			pic, err := w.cloud.DownLoadPic(item.PicKey)
			if err != nil {
				logrus.Errorf("Worker: download pictures from s3Cloud failed %v", err)

				errChan <- fmt.Errorf("worker: download pictures from s3Cloud failed %v", err)

				return
			}

			if err := w.sendPic(item.UserId, pic, item.Tag); err != nil {
				errChan <- fmt.Errorf("sendPic failed for %s: %w", item.PicKey, err)

				return
			}
		}(item)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		errs := getErrsFromChan(errChan)
		return errs
	}

	return nil
}

func getErrsFromChan(errChan chan error) []error {
	var errs []error

	for err := range errChan {
		errs = append(errs, err)
	}

	return errs
}

func (w *Worker) sendPic(chatID int64, pic []byte, tag string) error {
	pictures := tgbotapi.NewPhoto(chatID, tgbotapi.FileBytes{
		Name:  tag,
		Bytes: pic,
	})

	pictures.Caption = tag

	_, err := w.bot.TgAPI.Send(pictures)
	if err != nil {
		return err
	}

	return nil
}
