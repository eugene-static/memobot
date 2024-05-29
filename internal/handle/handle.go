package handle

import (
	"context"
	"log/slog"

	"github.com/eugene-static/memobot/internal/entities"
	"github.com/eugene-static/memobot/internal/service"
	"github.com/eugene-static/memobot/internal/session"
	"github.com/eugene-static/memobot/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	start     = "/start"
	addFolder = "/add_folder"
	addNote   = "/add_note"
	rename    = "/rename"
	update    = "/update"
	del       = "/delete"
	back      = "/back"
	accept    = "/accept"
	cancel    = "/cancel"
)
const (
	lvlRoot = iota
	lvlDir
	lvlNote
	lvlAction
	lvlAccept
)

const (
	folder = "📁"
	note   = "📄"
)

type Service interface {
	GetList(context.Context, int64, string) ([]*entities.List, error)
	Get(context.Context, int64, string) (string, error)
	AddRoot(context.Context, int64, string, string) error
	Add(context.Context, int64, string, string, bool) (string, error)
	UpdateContent(context.Context, string, string) error
	Rename(context.Context, string, string) error
	Delete(context.Context, string) error
}

type Handler struct {
	l       *slog.Logger
	service Service
	manager *session.Manager
}

func New(l *slog.Logger, service *service.Service, manager *session.Manager) *Handler {
	return &Handler{
		l:       l,
		service: service,
		manager: manager,
	}
}

func (h *Handler) Send(ctx context.Context, b *tgbotapi.BotAPI, u tgbotapi.Update) {
	var msg *tgbotapi.MessageConfig
	chat := u.FromChat()
	userID, username := chat.ID, chat.UserName
	user := h.manager.GetUser(userID)
	if user == nil {
		user = h.manager.AddUser(userID, username)
		if err := h.service.AddRoot(ctx, user.ID, user.CurrentDir().ID, user.Username); err != nil {
			h.l.Warn("failed to add root dir", slog.Int64("user_id", user.ID), logger.Err(err))
		}
	}
	if u.Message != nil {
		msg = h.Message(ctx, u.Message.Text, user)
	} else if u.CallbackQuery != nil {
		if user.LastMessageID != u.CallbackQuery.Message.MessageID {
			return
		}
		msg = h.Callback(ctx, u.CallbackQuery.Data, user)
	}
	if msg == nil {
		return
	}
	msg.ParseMode = "HTML"
	m, err := b.Send(*msg)
	user.LastMessageID = m.MessageID
	if err != nil {
		h.l.Warn("failed to send message", slog.Int64("chat_id", msg.ChatID), logger.Err(err))
	}
}

func errorMsg(chatID int64, msg string) *tgbotapi.MessageConfig {
	config := tgbotapi.NewMessage(chatID, msg)
	return &config
}

func getLevel(id string) int {
	if len(id) == 0 {
		return -1
	}
	switch id[0] {
	case '0':
		return lvlRoot
	case '1':
		return lvlDir
	case '2':
		return lvlNote
	}
	return -1
}
