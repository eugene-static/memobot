package handle

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/eugene-static/memobot/internal/session"
	"github.com/eugene-static/memobot/pkg/format"
	"github.com/eugene-static/memobot/pkg/logger"
	"github.com/eugene-static/memobot/pkg/wrapper"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) Callback(ctx context.Context, u tgbotapi.Update, user *session.User) *tgbotapi.MessageConfig {
	user.NewAction(u.CallbackQuery.Data)
	dir := user.CurrentDir()
	mc := &MessageConfig{
		userID: user.ID,
		title:  dir.Title,
		list:   nil,
	}
	h.l.Info("update", slog.String("callback", user.Action), slog.Int64("user_id", user.ID))
	l := h.l.With(
		slog.Int64("user_id", user.ID),
		slog.String("section_id", dir.ID),
		slog.String("request", user.Action),
	)
	switch user.Action {
	case addFolder:
		mc.msg = "Введите название раздела:"
		mc.level = lvlAction
	case addNote:
		mc.msg = "Введите название заметки:"
		mc.level = lvlAction
	case update:
		mc.msg = "Добавьте текст заметки:"
		mc.level = lvlAction
	case rename:
		mc.msg = fmt.Sprintf("Введите новое название для %s:", format.Format(dir.Title, format.Italic))
		mc.level = lvlAction
	case del:
		mc.msg = fmt.Sprintf("Удалить %s?", format.Format(dir.Title, format.Italic))
		mc.level = lvlAccept
	case accept:
		if err := h.service.Delete(ctx, dir.ID); err != nil {
			l.Warn("failed to delete section", logger.Err(err))
		}
		l.Info("deleted")
		fallthrough
	case back:
		user.Up()
		dir = user.CurrentDir()
		fallthrough
	case cancel:
		mc.level = getLevel(dir.ID)
		mc.msg = format.Format(user.Path(), format.Italic)
	default:
		mc.msg = format.Format(user.Down(wrapper.Unwrap(user.Action)), format.Italic)
		dir = user.CurrentDir()
		mc.level = getLevel(dir.ID)
	}
	if mc.level > lvlNote {
		return mc.build()
	}
	if mc.level == lvlNote {
		content, err := h.service.Get(ctx, mc.userID, dir.ID)
		if err != nil {
			l.Warn("failed to get content", logger.Err(err))
			return nil
		}
		mc.msg = fmt.Sprintf("%s\n\n%s", mc.msg, format.Format(content, format.Monotype))
		return mc.build()
	}
	list, err := h.service.GetList(ctx, user.ID, dir.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		l.Warn("failed to get list", logger.Err(err))
		return nil
	}
	mc.list = listToButton(list)
	return mc.build()
}
