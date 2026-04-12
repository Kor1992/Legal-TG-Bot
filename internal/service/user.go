package service

import (
	"context"
	"errors"
	"log"

	"github.com/Kor1992/Legal-TG-Bot/internal/domain"
	"github.com/Kor1992/Legal-TG-Bot/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Register(
	ctx context.Context,
	telegramID int64,
	firstName string, lastName string,
	phone string) (*domain.User, error) {
	user := &domain.User{
		TelegramID: telegramID,
		FirstName:  firstName,
		LastName:   &lastName,
		Phone:      &phone,
		Role:       "client",
	}

	err := s.repo.Create(ctx, user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetOrCreate(ctx context.Context,
	telegramID int64,
	firstName string) (*domain.User, error) {

	user, err := s.repo.GetByTelegramID(ctx, telegramID)
	if err == nil {
		return user, nil
	}

	if !errors.Is(err, repository.ErrNotFound) {
		log.Printf("Error getting user: %v", err)
		return nil, err
	}
	user = &domain.User{
		TelegramID: telegramID,
		FirstName:  firstName,
		LastName:   nil,
		Phone:      nil,
		Role:       "client",
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdatePhone(ctx context.Context, telegramID int64, phone string) error {
	user, err := s.repo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return err
	}

	user.Phone = &phone
	return s.repo.Update(ctx, user)
}

func (s *UserService) UpdateLastname(ctx context.Context, telegramID int64, lastName string) error {
	user, err := s.repo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return err
	}

	user.Phone = &lastName
	return s.repo.Update(ctx, user)
}
