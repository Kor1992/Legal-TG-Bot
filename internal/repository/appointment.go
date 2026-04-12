package repository

import (
	"context"
	"time"

	"github.com/Kor1992/Legal-TG-Bot/internal/domain"
)

type AppointmentRepository interface {
	Create(ctx context.Context, appointment *domain.Appointment) error
	GetByID(ctx context.Context, id int) (*domain.Appointment, error)
	ListByClient(ctx context.Context, clientID int) ([]*domain.Appointment, error)
	ListByLawyer(ctx context.Context, lawyerID int) ([]*domain.Appointment, error)
	ListByLawyerAndDate(ctx context.Context, lawyerID int, date time.Time) ([]*domain.Appointment, error)
	UpdateStatus(ctx context.Context, id int, status string) error
	Delete(ctx context.Context, id int) error
}
