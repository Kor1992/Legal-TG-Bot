package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/Kor1992/Legal-TG-Bot/internal/domain"
	"github.com/Kor1992/Legal-TG-Bot/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AppointmentRepo struct {
	db *pgxpool.Pool
}

func NewAppointmentRepository(db *pgxpool.Pool) repository.AppointmentRepository {
	return &AppointmentRepo{db: db}
}

func (r *AppointmentRepo) Create(ctx context.Context, appointment *domain.Appointment) error {
	query := `
		INSERT INTO ltb.appointments (client_id, lawyer_id, appointment_time, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at;
		`

	row := r.db.QueryRow(ctx, query,
		appointment.ClientID,
		appointment.LawyerID,
		appointment.AppointmentTime,
		appointment.Status,
	)

	return row.Scan(&appointment.ID, &appointment.CreatedAt)
}
func (r *AppointmentRepo) GetByID(ctx context.Context, id int) (*domain.Appointment, error) {
	query := `
		SELECT id, client_id, lawyer_id, appointment_time, status, created_at
		FROM ltb.appointments
		WHERE id = $1;
		`

	appointment := &domain.Appointment{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&appointment.ID,
		&appointment.ClientID,
		&appointment.LawyerID,
		&appointment.AppointmentTime,
		&appointment.Status,
		&appointment.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return appointment, nil
}
func (r *AppointmentRepo) ListByClient(ctx context.Context, clientID int) ([]*domain.Appointment, error) {
	query := `
		SELECT id, client_id, lawyer_id, appointment_time, status, created_at
		FROM ltb.appointments
		WHERE client_id = $1
		ORDER BY appointment_time DESC;
		`

	rows, err := r.db.Query(ctx, query, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appoinments []*domain.Appointment

	for rows.Next() {
		a := &domain.Appointment{}
		err := rows.Scan(
			&a.ID,
			&a.ClientID,
			&a.LawyerID,
			&a.AppointmentTime,
			&a.Status,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		appoinments = append(appoinments, a)
	}

	return appoinments, nil
}
func (r *AppointmentRepo) ListByLawyer(ctx context.Context, lawyerID int) ([]*domain.Appointment, error) {
	query := `
		SELECT id, client_id, lawyer_id, appointment_time, status, created_at
		FROM ltb.appointments
		WHERE lawyer_id = $1
		ORDER BY appointment_time DESC;
		`

	rows, err := r.db.Query(ctx, query, lawyerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*domain.Appointment
	for rows.Next() {
		a := &domain.Appointment{}
		err := rows.Scan(
			&a.ID,
			&a.ClientID,
			&a.LawyerID,
			&a.AppointmentTime,
			&a.Status,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, a)
	}

	return appointments, nil
}
func (r *AppointmentRepo) ListByLawyerAndDate(ctx context.Context, lawyerID int, date time.Time) ([]*domain.Appointment, error) {
	startOfDay := date.Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := `
		SELECT id, client_id, lawyer_id, appointment_time, status, created_at
		FROM ltb.appointments
		WHERE lawyer_id = $1 
		  AND appointment_time >= $2 
		  AND appointment_time < $3
		  AND status != 'cancelled'
		ORDER BY appointment_time`

	rows, err := r.db.Query(ctx, query, lawyerID, startOfDay, endOfDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*domain.Appointment
	for rows.Next() {
		a := &domain.Appointment{}
		err := rows.Scan(
			&a.ID,
			&a.ClientID,
			&a.LawyerID,
			&a.AppointmentTime,
			&a.Status,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, a)
	}

	return appointments, nil
}
func (r *AppointmentRepo) UpdateStatus(ctx context.Context, id int, status string) error {
	query := `UPDATE ltb.appointments SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, status, id)
	return err
}
func (r *AppointmentRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM ltb.appointments WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
