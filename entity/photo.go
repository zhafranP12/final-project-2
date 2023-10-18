package entity

import "time"

type photo struct {
	ID        int
	UserID    int
	Title     string
	Caption   string
	PhotoURL  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
