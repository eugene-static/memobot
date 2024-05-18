package server

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/eugene-static/memobot/internal/handle"
	"github.com/eugene-static/memobot/internal/service"
	"github.com/eugene-static/memobot/internal/session"
	"github.com/eugene-static/memobot/internal/storage/sqlite3"
	"github.com/eugene-static/memobot/pkg/config"
	"github.com/eugene-static/memobot/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Application struct {
	l   *slog.Logger
	bot *tgbotapi.BotAPI
	cfg *config.Config
}

func Init(l *slog.Logger, b *tgbotapi.BotAPI, cfg *config.Config) *Application {
	return &Application{
		l:   l,
		bot: b,
		cfg: cfg,
	}
}

func (a *Application) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	appStorage, err := storage.New(ctx, &a.cfg.Database, a.l)
	if err != nil {
		a.l.Error("database connection failed", logger.Err(err))
		return
	}
	a.l.Info("Authorized successful", slog.String("account", a.bot.Self.UserName))
	appService := service.New(appStorage)
	appSessionManager := session.New()
	appHandler := handle.NewHandler(a.l, appService, appSessionManager)
	a.bot.Debug = a.cfg.Bot.DebugMode
	u := tgbotapi.NewUpdate(a.cfg.Bot.UpdateOffset)
	u.Timeout = a.cfg.Bot.UpdateTimeout
	updates := a.bot.GetUpdatesChan(u)
	go func() {
		for update := range updates {
			appHandler.Send(ctx, a.bot, update)
		}
	}()
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	long := make(chan struct{}, 1)
	go func() {
		a.bot.StopReceivingUpdates()
		if err = appStorage.Close(); err != nil {
			a.l.Error("failed closing db", logger.Err(err))
		}
		long <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		a.l.Error("error during closing app", logger.Err(ctx.Err()))
	case <-long:
		a.l.Info("app shut down successfully")
	}
}
