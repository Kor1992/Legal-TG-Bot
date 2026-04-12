package domain

import "time"

type Lawyer struct {
	ID          int
	UserID      int
	FullName    string
	Description *string
	IsActive    bool
	CreatedAt   time.Time
}
