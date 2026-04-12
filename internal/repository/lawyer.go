package repository

import (
	"context"

	"github.com/Kor1992/Legal-TG-Bot/internal/domain"
)

type LawyerRepository interface {
	Create(ctx context.Context, lawyer *domain.Lawyer) error
	GetByUserID(ctx context.Context, userID int) (*domain.Lawyer, error)
	GetByID(ctx context.Context, id int) (*domain.Lawyer, error)
	ListActive(ctx context.Context) ([]*domain.Lawyer, error)
	Update(ctx context.Context, lawyer *domain.Lawyer) error
}
