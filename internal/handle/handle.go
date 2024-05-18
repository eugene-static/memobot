package handle

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/eugene-static/memobot/internal/entities"
	"github.com/eugene-static/memobot/internal/service"
	"github.com/eugene-static/memobot/internal/session"
	"github.com/eugene-static/memobot/pkg/logger"
	"github.com/eugene-static/memobot/pkg/wrapper"
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
	GetList(int64, string) ([]*entities.List, error)
	Get(int64, string) (string, error)
	AddRoot(int64, string, string) error
	Add(int64, string, string, bool) (string, error)
	UpdateContent(string, string) error
	Rename(string, string) error
	Delete(string) error
}

type Handler struct {
	l       *slog.Logger
	service Service
	manager *session.Manager
}

func NewHandler(l *slog.Logger, service *service.Service, manager *session.Manager) *Handler {
	return &Handler{
		l:       l,
		service: service,
		manager: manager,
	}
}

func (h *Handler) Send(b *tgbotapi.BotAPI, u tgbotapi.Update) {
	var msg *tgbotapi.MessageConfig
	var user *session.User
	if u.Message != nil {
		userID, username := u.Message.Chat.ID, u.Message.Chat.UserName
		user = h.manager.GetUser(userID, username)
		if user == nil {
			user = h.manager.AddUser(userID, username)
			if err := h.service.AddRoot(user.ID, user.CurrentDir().ID, user.Username); err != nil {
				h.l.Warn("failed to add root dir", slog.Int64("user_id", user.ID), logger.Err(err))
			}
		}
		msg = h.Message(u, user)
	} else if u.CallbackQuery != nil {
		userID, username := u.CallbackQuery.Message.Chat.ID, u.CallbackQuery.Message.Chat.UserName
		user = h.manager.GetUser(userID, username)
		if user == nil {
			user = h.manager.AddUser(userID, username)
			if err := h.service.AddRoot(user.ID, user.CurrentDir().ID, user.Username); err != nil {
				h.l.Warn("failed to add root dir", slog.Int64("user_id", user.ID), logger.Err(err))
			}
		}
		if user.LastMessageID != u.CallbackQuery.Message.MessageID {
			return
		}
		msg = h.Callback(u, user)
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

func listToButton(list []*entities.List) [][]tgbotapi.InlineKeyboardButton {
	buttons := make([][]tgbotapi.InlineKeyboardButton, len(list))
	symbol := folder
	for i, el := range list {
		if strings.HasPrefix(el.ID, "2") {
			symbol = note
		}
		title := fmt.Sprintf("%s%s", symbol, el.Title)
		id := wrapper.Wrap(el.ID, el.Title)
		buttons[i] = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(title, id))
	}
	return buttons
}
