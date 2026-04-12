package bot

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Kor1992/Legal-TG-Bot/internal/domain"
)

func PhoneKeyboard() tgbotapi.ReplyKeyboardMarkup {
	contactButton := tgbotapi.NewKeyboardButtonContact("📱 Отправить номер телефона")
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(contactButton),
	)
	keyboard.OneTimeKeyboard = true
	keyboard.ResizeKeyboard = true
	return keyboard
}

func RemoveKeyboard() tgbotapi.ReplyKeyboardRemove {
	return tgbotapi.NewRemoveKeyboard(true)
}

func LawyersKeyboard(lawyers []*domain.Lawyer) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, lawyer := range lawyers {
		text := lawyer.FullName
		if lawyer.Description != nil && *lawyer.Description != "" {
			text += " — " + *lawyer.Description
		}
		button := tgbotapi.NewInlineKeyboardButtonData(text, "lawyer_"+strconv.Itoa(lawyer.ID))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func SlotsKeyboard(slots []Slot) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, slot := range slots {
		button := tgbotapi.NewInlineKeyboardButtonData(slot.Label, slot.Callback)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
