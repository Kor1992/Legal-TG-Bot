package domain

import "time"

type Appointment struct {
	ID              int
	ClientID        int
	LawyerID        int
	AppointmentTime time.Time
	Status          string
	CreatedAt       time.Time
}
