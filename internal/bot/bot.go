package bot

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Kor1992/Legal-TG-Bot/internal/service"
)

type Bot struct {
	api                *tgbotapi.BotAPI
	userService        *service.UserService
	lawyerService      *service.LawyerService
	appointmentService *service.AppointmentService
	stateManager       *StateManager
}

func NewBot(
	token string,
	userService *service.UserService,
	lawyerService *service.LawyerService,
	appointmentService *service.AppointmentService,
) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	api.Debug = true
	log.Printf("Authorized on account %s", api.Self.UserName)

	return &Bot{
		api:                api,
		userService:        userService,
		lawyerService:      lawyerService,
		appointmentService: appointmentService,
		stateManager:       NewStateManager(),
	}, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		b.handleUpdate(update)
	}
}

func (b *Bot) handleUpdate(update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		b.handleCallback(update.CallbackQuery)
		return
	}

	// Обработка сообщений
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	telegramID := update.Message.From.ID
	text := update.Message.Text

	log.Printf("[%s] %s", update.Message.From.UserName, text)

	state := b.stateManager.Get(telegramID)

	switch state {
	case StateAwaitingName:
		b.handleNameInput(chatID, telegramID, text)
	case StateAwaitingPhone:
		b.handlePhoneInput(chatID, telegramID, update.Message)
	default:
		if strings.HasPrefix(text, "/") {
			b.handleCommand(chatID, telegramID, text)
		} else {
			b.handleDefault(chatID, telegramID, text, update.Message.From.FirstName)
		}
	}
}

func (b *Bot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	b.api.Send(msg)
}

func (b *Bot) sendMessageWithKeyboard(chatID int64, text string, keyboard interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	b.api.Send(msg)
}
