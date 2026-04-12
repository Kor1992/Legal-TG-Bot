package postgres

import (
	"context"
	"errors"

	"github.com/Kor1992/Legal-TG-Bot/internal/domain"
	"github.com/Kor1992/Legal-TG-Bot/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LawyerRepo struct {
	db *pgxpool.Pool
}

func NewLawyerRepository(db *pgxpool.Pool) repository.LawyerRepository {
	return &LawyerRepo{db: db}
}
func (r *LawyerRepo) Create(ctx context.Context, lawyer *domain.Lawyer) error {
	query := `
	INSERT INTO  ltb.lawyers (user_id, full_name, description, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at;
	`
	row := r.db.QueryRow(ctx, query,
		lawyer.UserID,
		lawyer.FullName,
		lawyer.Description,
		lawyer.IsActive)
	return row.Scan(&lawyer.ID, lawyer.CreatedAt)

}
func (r *LawyerRepo) GetByUserID(ctx context.Context, userID int) (*domain.Lawyer, error) {
	query := `
	SELECT id, user_id, full_name, description, is_active, created_at
		FROM ltb.lawyers
		WHERE user_id = $1;
	`
	lawyer := &domain.Lawyer{}

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&lawyer.ID,
		&lawyer.UserID,
		&lawyer.FullName,
		&lawyer.Description,
		&lawyer.IsActive,
		&lawyer.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return lawyer, nil
}
func (r *LawyerRepo) GetByID(ctx context.Context, id int) (*domain.Lawyer, error) {
	query := `
		SELECT id, user_id, full_name, description, is_active, created_at
		FROM ltb.lawyers
		WHERE id = $1`

	lawyer := &domain.Lawyer{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&lawyer.ID,
		&lawyer.UserID,
		&lawyer.FullName,
		&lawyer.Description,
		&lawyer.IsActive,
		&lawyer.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return lawyer, nil
}

func (r *LawyerRepo) ListActive(ctx context.Context) ([]*domain.Lawyer, error) {
	query := `
	SELECT id, user_id, full_name, description, is_active, created_at
	FROM ltb.lawyers
	WHERE is_active = true
	ORDER BY full_name;
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lawyers []*domain.Lawyer

	for rows.Next() {
		lawyer := &domain.Lawyer{}
		err := rows.Scan(
			&lawyer.ID,
			&lawyer.UserID,
			&lawyer.FullName,
			&lawyer.Description,
			&lawyer.IsActive,
			&lawyer.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		lawyers = append(lawyers, lawyer)

	}
	return lawyers, nil

}

func (r *LawyerRepo) Update(ctx context.Context, lawyer *domain.Lawyer) error {
	query := `
	UPDATE ltb.lawyers
	SET full_name = $1, description = $2, is_active = $3
	WHERE id = $4;
	`

	_, err := r.db.Exec(ctx, query, lawyer.FullName,
		lawyer.Description,
		lawyer.IsActive,
		lawyer.ID)

	return err
}
