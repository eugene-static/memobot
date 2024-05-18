package handle

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageConfig struct {
	userID int64
	title  string
	level  int
	msg    string
	list   [][]tgbotapi.InlineKeyboardButton
}

func (mc *MessageConfig) build() *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(mc.userID, mc.msg)
	switch mc.level {
	case lvlRoot:
		mc.list = append(mc.list,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Добавить раздел", addFolder),
				tgbotapi.NewInlineKeyboardButtonData("Создать запись", addNote)))
	case lvlDir:
		mc.list = append(mc.list,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Добавить раздел", addFolder),
				tgbotapi.NewInlineKeyboardButtonData("Создать запись", addNote)),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Переименовать", rename),
				tgbotapi.NewInlineKeyboardButtonData("Удалить", del)),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Назад", back)))
	case lvlNote:
		mc.list = append(mc.list,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Переименовать", rename),
				tgbotapi.NewInlineKeyboardButtonData("Редактировать", update)),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Удалить", del)),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Назад", back)))
	case lvlAction:
		mc.list = append(mc.list,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Отмена", cancel)))
	case lvlAccept:
		mc.list = append(mc.list,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("OK", accept),
				tgbotapi.NewInlineKeyboardButtonData("Отмена", cancel)))
	}
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(mc.list...)
	return &msg

}
