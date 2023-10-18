package dto

import "time"

type NewSocialMediaRequest struct {
	Name           string `json:"name"`
	SocialMediaURL string `json:"social_media_url"`
	UserID         int    `json:"user_id"`
}

type NewSocialMediaResponse struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	SocialMediaURL string    `json:"social_media_url"`
	UserID         int       `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
}

type GetUsersSocialMedia struct {
	ID              int    `json:"id"`
	Username        string `json:"username"`
	ProfileImageURL string `json:"profil_image_url"`
}

type GetSocialMediaWithUser struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	SocialMediaURL string    `json:"social_media_url"`
	UserID         int       `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	User           GetUsersSocialMedia
}

type GetSocialMediaResponse struct {
	Data []GetSocialMediaWithUser
}

type UpdateSocialMediaRequest struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	SocialMediaURL string `json:"social_media_url"`
	UserID         int    `json:"user_id"`
}

type UpdateSocialMediaResponse struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	SocialMediaURL string    `json:"social_media_url"`
	UserID         int       `json:"user_id"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type DeleteSocialMediaResponse struct {
	Message string `json:"message"`
}
