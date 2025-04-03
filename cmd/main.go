package main

import (
	tgbot "MemoryPicBot/bot"
	"MemoryPicBot/bot/config"
	"MemoryPicBot/bot/handlers"
	"MemoryPicBot/bot/handlers/state"
	"MemoryPicBot/bot/repository"
	worker2 "MemoryPicBot/bot/worker"
	"MemoryPicBot/pkg/cache"
	"MemoryPicBot/pkg/database"
	"MemoryPicBot/pkg/s3"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg, err := config.InitConfig()
	if err != nil {
		logrus.Panicf("cant init config: %v ", err)
	}

	db, err := database.InitDB(&cfg.DBConfig)
	if err != nil {
		logrus.Fatalf("cant init db: %v ", err)
	}
	defer db.Close()

	redisClient, err := cache.NewRedisClient(cfg.RedisConfig)
	if err != nil {
		logrus.Panicf("cant init Redis: %v", err)
	}
	defer redisClient.Close()

	s3Client, err := s3.NewS3Client(cfg.AWSConfig.Region, cfg.AWSConfig.Bucket)
	if err != nil {
		log.Fatalf("cant init s3 client: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		logrus.Fatalf("invalid telegram token: %v", err)
	}

	bot.Debug = true

	repo := repository.NewPicRepository(db)

	userState := state.NewStateManager()
	tgBot := tgbot.NewBot(bot, cfg, userState)
	h := handlers.InitHandlers(tgBot, redisClient, repo, s3Client, userState)
	worker := worker2.NewWorker(repo, tgBot, s3Client)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)

	logrus.Info("Starting bot...")

	go func() {
		defer wg.Done()
		tgBot.Start(h)
		logrus.Info("Bot stopped")
	}()

	wg.Add(1)

	go func() {
		logrus.Info("Worker goroutine started")
		defer wg.Done()
		if err := worker.Run(ctx); err != nil {
			logrus.Errorf("Worker stopped with error: %v", err)
		}
	}()

	<-signals
	logrus.Warn("Received shutdown signal, stopping services...")

	cancel()

	wg.Wait()
}
