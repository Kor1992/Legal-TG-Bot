package service

import (
	"context"
	"time"

	"github.com/Kor1992/Legal-TG-Bot/internal/domain"
	"github.com/Kor1992/Legal-TG-Bot/internal/repository"
)

type AppointmentService struct {
	appointmentRepo repository.AppointmentRepository
}

func NewAppointmentService(appointmentRepo repository.AppointmentRepository) *AppointmentService {
	return &AppointmentService{
		appointmentRepo: appointmentRepo,
	}
}

func (s *AppointmentService) Create(ctx context.Context, clientID, lawyerID int, appointmentTime time.Time) (*domain.Appointment, error) {
	appointment := &domain.Appointment{
		ClientID:        clientID,
		LawyerID:        lawyerID,
		AppointmentTime: appointmentTime,
		Status:          "confirmed",
	}

	err := s.appointmentRepo.Create(ctx, appointment)
	if err != nil {
		return nil, err
	}

	return appointment, nil
}

func (s *AppointmentService) ListByClient(ctx context.Context, clientID int) ([]*domain.Appointment, error) {
	return s.appointmentRepo.ListByClient(ctx, clientID)
}
