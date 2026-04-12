package main

import (
	"context"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/Kor1992/Legal-TG-Bot/internal/repository/postgres"
	"github.com/Kor1992/Legal-TG-Bot/internal/service"
)

const (
	stateAwaitingName  = "awaiting_name"
	stateAwaitingPhone = "awaiting_phone"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set in .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")

		if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
			log.Fatal("Database configuration is incomplete")
		}

		dbURL = "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName
	}

	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}
	log.Println("Connected to database")

	userRepo := postgres.NewUserRepository(dbPool)
	userService := service.NewUserService(userRepo)

	userStates := make(map[int64]string)
	tempNames := make(map[int64]string)

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		telegramID := update.Message.From.ID
		text := update.Message.Text

		log.Printf("[%s] %s", update.Message.From.UserName, text)

		state := userStates[telegramID]

		switch state {
		case stateAwaitingName:
			handleNameInput(bot, chatID, telegramID, text, userStates, tempNames)

		case stateAwaitingPhone:
			handlePhoneInput(bot, chatID, telegramID, update.Message, userStates, tempNames, userService)

		default:
			handleDefault(bot, chatID, telegramID, text, userStates, tempNames, userService, update.Message.From.FirstName)
		}
	}
}

func handleDefault(
	bot *tgbotapi.BotAPI,
	chatID int64,
	telegramID int64,
	text string,
	userStates map[int64]string,
	tempNames map[int64]string,
	userService *service.UserService,
	firstName string,
) {
	switch text {
	case "/start":
		ctx := context.Background()
		user, err := userService.GetOrCreate(ctx, telegramID, "")

		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Произошла ошибка. Попробуйте позже.")
			bot.Send(msg)
			log.Printf("Error in GetOrCreate: %v", err)
			return
		}

		if user.Phone == nil || *user.Phone == "" {
			userStates[telegramID] = stateAwaitingName

			msg := tgbotapi.NewMessage(chatID, "Добро пожаловать! Для записи на консультацию мне нужно знать ваше имя.\n\nКак вас зовут?")
			bot.Send(msg)
		} else {
			displayName := user.FirstName
			if displayName == "" {
				displayName = "уважаемый пользователь"
			}
			msg := tgbotapi.NewMessage(chatID, "С возвращением, "+displayName+"! 👋\n\nВы уже зарегистрированы. Используйте меню для навигации.")
			bot.Send(msg)
		}

	default:
		ctx := context.Background()
		user, err := userService.GetOrCreate(ctx, telegramID, firstName)

		if err != nil {
			log.Printf("Error in GetOrCreate: %v", err)
			return
		}

		if user.Phone == nil || *user.Phone == "" {
			userStates[telegramID] = stateAwaitingName
			msg := tgbotapi.NewMessage(chatID, "Давайте познакомимся! Как вас зовут?")
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(chatID, "Вы сказали: "+text)
			bot.Send(msg)
		}
	}
}

func handleNameInput(
	bot *tgbotapi.BotAPI,
	chatID int64,
	telegramID int64,
	name string,
	userStates map[int64]string,
	tempNames map[int64]string,
) {
	tempNames[telegramID] = name

	userStates[telegramID] = stateAwaitingPhone

	contactButton := tgbotapi.NewKeyboardButtonContact("📱 Отправить номер телефона")
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(contactButton),
	)
	keyboard.OneTimeKeyboard = true
	keyboard.ResizeKeyboard = true

	msg := tgbotapi.NewMessage(chatID, "Спасибо, "+name+"! Теперь отправьте ваш номер телефона.\n\nНажмите кнопку ниже, чтобы поделиться контактом:")
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func handlePhoneInput(
	bot *tgbotapi.BotAPI,
	chatID int64,
	telegramID int64,
	message *tgbotapi.Message,
	userStates map[int64]string,
	tempNames map[int64]string,
	userService *service.UserService,
) {
	var phone string

	if message.Contact != nil {
		phone = message.Contact.PhoneNumber
	} else {
		phone = message.Text
	}

	name := tempNames[telegramID]
	delete(tempNames, telegramID)

	delete(userStates, telegramID)

	ctx := context.Background()

	user, err := userService.GetOrCreate(ctx, telegramID, name)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Произошла ошибка при сохранении данных. Попробуйте позже.")
		bot.Send(msg)
		log.Printf("Error in GetOrCreate: %v", err)
		return
	}

	if name != "" {
		user.FirstName = name
	}
	if phone != "" {
		user.Phone = &phone
	}

	err = userService.UpdatePhone(ctx, telegramID, phone)
	if err != nil {
		log.Printf("Error updating user: %v", err)
	}

	msg := tgbotapi.NewMessage(chatID, "✅ Регистрация успешно завершена!\n\nТеперь вы можете:\n• Записаться на консультацию\n• Посмотреть свои записи\n• Получить напоминание перед консультацией")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	bot.Send(msg)
}
