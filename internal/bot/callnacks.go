package bot

import (
	"context"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleCallback(callback *tgbotapi.CallbackQuery) {
	data := callback.Data
	chatID := callback.Message.Chat.ID
	telegramID := callback.From.ID

	b.api.Request(tgbotapi.NewCallback(callback.ID, ""))

	ctx := context.Background()
	user, err := b.userService.GetOrCreate(ctx, telegramID, "")
	if err != nil || user.Phone == nil || *user.Phone == "" {
		b.sendMessage(chatID, "Сначала зарегистрируйтесь: /start")
		return
	}

	if strings.HasPrefix(data, "lawyer_") {
		lawyerID, _ := strconv.Atoi(strings.TrimPrefix(data, "lawyer_"))
		b.showAvailableSlots(chatID, lawyerID)
	} else if strings.HasPrefix(data, "slot_") {
		lawyerID, slotTime, ok := ParseCallback(data)
		if ok {
			b.bookAppointment(chatID, user.ID, lawyerID, slotTime)
		}
	}
}

func (b *Bot) showAvailableSlots(chatID int64, lawyerID int) {
	slots := GetAvailableSlots(lawyerID, 2) // на 2 дня вперёд

	if len(slots) == 0 {
		b.sendMessage(chatID, "Нет доступных слотов")
		return
	}

	b.sendMessageWithKeyboard(chatID, "Выберите удобное время:", SlotsKeyboard(slots))
}

func (b *Bot) bookAppointment(chatID int64, clientID int, lawyerID int, appointmentTime time.Time) {
	ctx := context.Background()
	_, err := b.appointmentService.Create(ctx, clientID, lawyerID, appointmentTime)
	if err != nil {
		b.sendMessage(chatID, "Ошибка при создании записи. Возможно, это время уже занято.")
		return
	}

	b.sendMessage(chatID, "✅ Вы записаны на консультацию!\n\nДата и время: "+appointmentTime.Format("02.01.2006 15:04"))
}
