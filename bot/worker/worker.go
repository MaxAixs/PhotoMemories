package worker

import (
	"MemoryPicBot/bot"
	"MemoryPicBot/bot/repository"
	"MemoryPicBot/pkg/s3"
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

type Worker struct {
	bot   *bot.Bot
	repo  *repository.PicRepository
	cloud *s3.Client
}

func NewWorker(repo *repository.PicRepository, bot *bot.Bot, cloud *s3.Client) *Worker {
	return &Worker{
		bot:   bot,
		repo:  repo,
		cloud: cloud,
	}
}

type PicGetterList interface {
	GetPicList() ([]bot.UserPictures, error)
}

func (w *Worker) Run(ctx context.Context) error {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logrus.Info("Worker is shutting down")

			return ctx.Err()

		case <-ticker.C:
			logrus.Info("Start getting pickLists")

			picLists, err := w.repo.GetPicList()
			if err != nil {
				logrus.Errorf("GetPicList error: %v", err)

				continue
			}

			if len(picLists) == 0 {
				logrus.Info("No pictures found")

				continue
			}

			if err := w.ProcessPicLists(picLists); err != nil {
				logrus.Errorf("error: %v", err)

				continue
			}
		}
	}
}
