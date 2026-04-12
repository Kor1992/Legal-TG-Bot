package repository

import (
	"context"

	"github.com/Kor1992/Legal-TG-Bot/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}
