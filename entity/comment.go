package entity

import "time"

type comment struct {
	ID        int
	UserID    int
	PhotoID   int
	Message   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
