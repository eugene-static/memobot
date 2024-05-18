package main

import (
	"log/slog"
	"os"

	"github.com/eugene-static/memobot/pkg/config"
	"github.com/eugene-static/memobot/server"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	sections      = 1
	subsections   = 2
	note          = 3
	sectionAdd    = 11
	subsectionAdd = 21
)

func main() {
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	}))
	cfg, err := config.GetConfig("./config/config.json")
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
