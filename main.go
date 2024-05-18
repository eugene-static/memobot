package main

import (
	"log/slog"
	"os"

	"github.com/eugene-static/memobot/internal/server"
	"github.com/eugene-static/memobot/pkg/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	}))
	cfg, err := config.Get("./config/config.json")
	if err != nil || cfg == nil {
		l.Error("failed to get config", slog.Any("details", err))
		return
	}
	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		l.Error("bot connection failed", slog.Any("details", err))
		return
	}
	app := server.Init(l, bot, cfg)
	app.Run()
}
