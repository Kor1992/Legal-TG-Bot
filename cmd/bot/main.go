package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/Kor1992/Legal-TG-Bot/internal/bot"
	"github.com/Kor1992/Legal-TG-Bot/internal/repository/postgres"
	"github.com/Kor1992/Legal-TG-Bot/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
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
	lawyerRepo := postgres.NewLawyerRepository(dbPool)
	appointmentRepo := postgres.NewAppointmentRepository(dbPool)

	userService := service.NewUserService(userRepo)
	lawyerService := service.NewLawyerService(lawyerRepo, userRepo)
	appointmentService := service.NewAppointmentService(appointmentRepo)

	b, err := bot.NewBot(token, userService, lawyerService, appointmentService)
	if err != nil {
		log.Fatal(err)
	}

	b.Start()
}
