package dto

import "time"

type NewUserRequest struct {
	Age      int    `json:"age" validate:"required,gt=8"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Username string `json:"username" validate:"required"`
}

type NewUserResponse struct {
	// status 201

	Age      int    `json:"age"`
	Email    string `json:"email"`
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type LoginRequest struct {
	// status 200
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type UpdateUserRequest struct {
	ID       int    `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
}

type UpdateUserResponse struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Age       int       `json:"age"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DeleteUserResponse struct {
	Message string `json:"message"`
}
