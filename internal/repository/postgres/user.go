package postgres

import (
	"context"
	"errors"

	"github.com/Kor1992/Legal-TG-Bot/internal/domain"
	"github.com/Kor1992/Legal-TG-Bot/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) repository.UserRepository {
	return &UserRepo{
		db: pool,
	}
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	query := `
	INSERT INTO ltb.users (telegram_id, first_name, last_name, phone, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at;
	`
	row := r.db.QueryRow(ctx,
		query,
		user.TelegramID,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.Role,
	)

	err := row.Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}

	return nil

}
func (r *UserRepo) GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	query := `
	SELECT id, telegram_id, first_name, last_name, phone, role, created_at
	FROM ltb.users
	WHERE telegram_id = $1;
`
	user := &domain.User{}

	err := r.db.QueryRow(ctx, query, telegramID).Scan(
		&user.ID,
		&user.TelegramID,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return user, nil

}
func (r *UserRepo) Update(ctx context.Context, user *domain.User) error {
	query := `
	UPDATE ltb.users
	SET first_name = $1, last_name = $2, phone = $3, role = $4
	WHERE id = $5;
`

	_, err := r.db.Exec(ctx, query, user.FirstName,
		user.LastName,
		user.Phone,
		user.Role,
		user.ID)

	return err
}
