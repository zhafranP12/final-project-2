package dto

import "time"

type NewCommentRequest struct {
	UserID  int    `json:"user_id"`
	Message string `json:"message" validate:"required"`
	PhotoID int    `json:"photo_id"`
}

type NewCommentResponse struct {
	ID        int       `json:"id"`
	Message   string    `json:"message"`
	PhotoID   int       `json:"photo_id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type GetPhotosComment struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoURL string `json:"photo_url"`
	UserID   int    `json:"user_id"`
}

type GetUsersComment struct {
	Email    string `json:"email"`
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type GetCommentsWithUserAndPhoto struct {
	ID        int       `json:"id"`
	Message   string    `json:"message"`
	PhotoID   int       `json:"photo_id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      GetUsersComment
	Photo     GetPhotosComment
}

type GetCommentsResponse struct {
	Data []GetCommentsWithUserAndPhoto
}

type UpdateCommentRequest struct {
	ID      int    `json:"id"`
	Message string `json:"message" validate:"required"`
	UserID  int    `json:"user_id"`
}

type UpdateCommentResponse struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoURL  string    `json:"photo_url"`
	UserID    int       `json:"user_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DeleteCommentResponse struct {
	Message string `json:"message"`
}

type DeleteCommentRequest struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`
}
