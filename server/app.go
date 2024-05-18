package server

import (
	"log/slog"

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
	b   *tgbotapi.BotAPI
	cfg *config.Config
}

func Init(l *slog.Logger, b *tgbotapi.BotAPI, cfg *config.Config) *Application {
	return &Application{
		l:   l,
		b:   b,
		cfg: cfg,
	}
}

func (a *Application) Run() {
	appStorage, err := storage.New(&a.cfg.Database, a.l)
	if err != nil {
		a.l.Error("database connection failed", logger.Err(err))
		return
	}
	a.l.Info("Authorized successful", slog.String("account", a.b.Self.UserName))
	appService := service.New(appStorage)
	appSessionManager := session.New()
	appHandler := handle.NewHandler(a.l, appService, appSessionManager)
	a.b.Debug = a.cfg.Bot.DebugMode
	u := tgbotapi.NewUpdate(a.cfg.Bot.UpdateOffset)
	u.Timeout = a.cfg.Bot.UpdateTimeout
	updates := a.b.GetUpdatesChan(u)
	for update := range updates {
		appHandler.Send(a.b, update)
	}
}
