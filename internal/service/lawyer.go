package service

import (
	"context"

	"github.com/Kor1992/Legal-TG-Bot/internal/domain"
	"github.com/Kor1992/Legal-TG-Bot/internal/repository"
)

type LawyerService struct {
	lawyerRepo repository.LawyerRepository
	userRepo   repository.UserRepository
}

func NewLawyerService(lawyerRepo repository.LawyerRepository, userRepo repository.UserRepository) *LawyerService {
	return &LawyerService{
		lawyerRepo: lawyerRepo,
		userRepo:   userRepo,
	}
}

func (s *LawyerService) ListActive(ctx context.Context) ([]*domain.Lawyer, error) {
	return s.lawyerRepo.ListActive(ctx)
}

func (s *LawyerService) GetByID(ctx context.Context, id int) (*domain.Lawyer, error) {
	return s.lawyerRepo.GetByID(ctx, id)
}
