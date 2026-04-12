package bot

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleDefault(chatID int64, telegramID int64, text string, firstName string) {
	ctx := context.Background()
	user, err := b.userService.GetOrCreate(ctx, telegramID, firstName)
	if err != nil {
		log.Printf("Error in GetOrCreate: %v", err)
		return
	}

	if user.Phone == nil || *user.Phone == "" {
		b.stateManager.Set(telegramID, StateAwaitingName)
		b.sendMessage(chatID, "Давайте познакомимся! Как вас зовут?")
	} else {
		b.sendMessage(chatID, "Используйте команды:\n/lawyers — список юристов\n/appointments — мои записи")
	}
}

func (b *Bot) handleNameInput(chatID int64, telegramID int64, name string) {
	b.stateManager.SetTempName(telegramID, name)
	b.stateManager.Set(telegramID, StateAwaitingPhone)

	b.sendMessageWithKeyboard(chatID,
		"Спасибо, "+name+"! Теперь отправьте ваш номер телефона.\n\nНажмите кнопку ниже, чтобы поделиться контактом:",
		PhoneKeyboard(),
	)
}

func (b *Bot) handlePhoneInput(chatID int64, telegramID int64, message *tgbotapi.Message) {
	var phone string
	if message.Contact != nil {
		phone = message.Contact.PhoneNumber
	} else {
		phone = message.Text
	}

	name := b.stateManager.GetTempName(telegramID)
	b.stateManager.Clear(telegramID)

	ctx := context.Background()
	user, err := b.userService.GetOrCreate(ctx, telegramID, name)
	if err != nil {
		b.sendMessage(chatID, "Ошибка при сохранении данных.")
		return
	}

	if name != "" {
		user.FirstName = name
	}
	if phone != "" {
		user.Phone = &phone
		b.userService.UpdatePhone(ctx, telegramID, phone)
	}

	b.sendMessageWithKeyboard(chatID,
		"✅ Регистрация успешно завершена!\n\nТеперь вы можете:\n/lawyers — список юристов\n/appointments — мои записи",
		RemoveKeyboard(),
	)
}

func (b *Bot) handleCommand(chatID int64, telegramID int64, text string) {
	ctx := context.Background()
	user, err := b.userService.GetOrCreate(ctx, telegramID, "")
	if err != nil {
		b.sendMessage(chatID, "Ошибка. Попробуйте позже.")
		return
	}

	switch text {
	case "/start":
		if user.Phone == nil || *user.Phone == "" {
			b.stateManager.Set(telegramID, StateAwaitingName)
			b.sendMessage(chatID, "Добро пожаловать! Для записи на консультацию мне нужно знать ваше имя.\n\nКак вас зовут?")
		} else {
			displayName := user.FirstName
			if displayName == "" {
				displayName = "уважаемый пользователь"
			}
			b.sendMessage(chatID, "С возвращением, "+displayName+"! 👋\n\nДоступные команды:\n/lawyers — список юристов\n/appointments — мои записи")
		}

	case "/lawyers":
		if user.Phone == nil || *user.Phone == "" {
			b.sendMessage(chatID, "Сначала зарегистрируйтесь: /start")
			return
		}
		b.showLawyersList(chatID)

	case "/appointments":
		if user.Phone == nil || *user.Phone == "" {
			b.sendMessage(chatID, "Сначала зарегистрируйтесь: /start")
			return
		}
		b.showMyAppointments(chatID, user.ID)

	default:
		b.sendMessage(chatID, "Неизвестная команда.\n/lawyers — список юристов\n/appointments — мои записи")
	}
}

func (b *Bot) showLawyersList(chatID int64) {
	ctx := context.Background()
	lawyers, err := b.lawyerService.ListActive(ctx)
	if err != nil {
		b.sendMessage(chatID, "Ошибка при загрузке списка юристов")
		return
	}

	if len(lawyers) == 0 {
		b.sendMessage(chatID, "Нет доступных юристов")
		return
	}

	b.sendMessageWithKeyboard(chatID, "Выберите юриста для консультации:", LawyersKeyboard(lawyers))
}

func (b *Bot) showMyAppointments(chatID int64, clientID int) {
	ctx := context.Background()
	appointments, err := b.appointmentService.ListByClient(ctx, clientID)
	if err != nil {
		b.sendMessage(chatID, "Ошибка при загрузке записей")
		return
	}

	if len(appointments) == 0 {
		b.sendMessage(chatID, "У вас пока нет записей.\nИспользуйте /lawyers чтобы записаться.")
		return
	}

	var text string
	for _, a := range appointments {
		timeStr := a.AppointmentTime.Format("02.01.2006 15:04")
		status := "✅"
		if a.Status == "cancelled" {
			status = "❌"
		}
		text += status + " " + timeStr + "\n"
	}

	b.sendMessage(chatID, "Ваши записи:\n"+text)
}
