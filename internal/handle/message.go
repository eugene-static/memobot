package handle

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/eugene-static/memobot/internal/session"
	"github.com/eugene-static/memobot/pkg/format"
	"github.com/eugene-static/memobot/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) Message(ctx context.Context, data string, user *session.User) *tgbotapi.MessageConfig {
	dir := user.CurrentDir()
	h.l.Info("update", slog.String("message", data), slog.Int64("user_id", user.ID))
	l := h.l.With(
		slog.Int64("user_id", user.ID),
		slog.String("element_id", dir.ID),
		slog.String("message", data),
	)
	if data == start {
		user = user.Root(user.Username)
		dir = user.CurrentDir()
	} else {
		switch user.Action {
		case addFolder:
			if strings.ContainsRune(data, '/') {
				return errorMsg(user.ID, "Название не должно содержать '/'")
			}
			id, err := h.service.Add(ctx, user.ID, dir.ID, data, true)
			if err != nil {
				l.Warn("failed to add new section", logger.Err(err))
				return nil
			}
			l.Info("section added", slog.String("new_section_id", id))
		case addNote:
			if strings.ContainsRune(data, '/') {
				return errorMsg(user.ID, "Название не должно содержать '/'")
			}
			id, err := h.service.Add(ctx, user.ID, dir.ID, data, false)
			if err != nil {
				l.Warn("failed to add new note", logger.Err(err))
				return nil
			}
			l.Info("note added", slog.String("new_note_id", id))
		case rename:
			if strings.ContainsRune(data, '/') {
				return errorMsg(user.ID, "Название не должно содержать '/'")
			}
			if err := h.service.Rename(ctx, dir.ID, data); err != nil {
				l.Warn("failed to rename", logger.Err(err))
				return nil
			}
			dir.Title = data
			l.Info("element renamed")
		case update:
			if err := h.service.UpdateContent(ctx, dir.ID, data); err != nil {
				l.Warn("failed to add content", logger.Err(err))
				return nil
			}
			l.Info("note updated")
		}
	}
	mc := &MessageConfig{
		userID: user.ID,
		title:  dir.Title,
		level:  getLevel(dir.ID),
		msg:    format.Format(user.Path(), format.Italic),
		list:   nil,
	}
	if mc.level == lvlNote {
		content, err := h.service.Get(ctx, user.ID, dir.ID)
		if err != nil {
			l.Warn("failed to get content")
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
	mc.listToButtons(list)
	return mc.build()
}
