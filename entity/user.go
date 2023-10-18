package entity

import "time"

type User struct {
	ID        int
	Email     string
	Username  string
	Password  string
	Age       int
	CreatedAt time.Time
	UpdatedAt time.Time
}
