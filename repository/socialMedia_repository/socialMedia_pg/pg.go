package social_media_pg

import (
	"database/sql"
	"finalProject2/dto"
	"finalProject2/pkg/errs"
	social_media_repository "finalProject2/repository/socialMedia_repository"
)

const (
	createSocialMedia = `
		INSERT INTO social_medias (name,social_media_url, user_id)
		VALUES ($1,$2,$3)
		RETURNING id,name,social_media_url,user_id,created_at
	`

	getSocialMedias = `
		SELECT 
			s.id,s.name,s.social_media_url,s.user_id,
			s.created_at,s.updated_at,
			u.id,u.username
		FROM social_medias AS s
		LEFT JOIN users AS u
			ON s.user_id = u.id
			WHERE s.user_id = $1
	`

	getUserId = `
		SELECT user_id FROM social_medias WHERE id = $1
	`

	updateSocialMedia = `
		UPDATE 
		social_medias SET name = $1, social_media_url = $2,
		updated_at = current_timestamp
		WHERE id = $3
		RETURNING id,name,social_media_url,user_id,updated_at
	`

	deleteSocialMedia = `
		DELETE FROM social_medias where id = $1
	`
)

type socialMediaPG struct {
	db *sql.DB
}

func NewSocialMediaPG(db *sql.DB) social_media_repository.Repository {
	return &socialMediaPG{
		db: db,
	}
}

func (socialMediaPG *socialMediaPG) CreateSocialMedia(sm dto.NewSocialMediaRequest) (*dto.NewSocialMediaResponse, errs.Error) {
	// request :  name,social_media_url, user_id
	// response : id,name,social_media_url,user_id,created_at
	var resp dto.NewSocialMediaResponse

	err := socialMediaPG.db.QueryRow(createSocialMedia, sm.Name, sm.SocialMediaURL, sm.UserID).Scan(
		&resp.ID, &resp.Name, &resp.SocialMediaURL, &resp.UserID, &resp.CreatedAt,
	)

	if err != nil {
		return nil, errs.NewInternalServerError(err.Error())
	}

	return &resp, nil
}

func (socialMediaPG *socialMediaPG) GetSocialMedias(userId int) (*dto.GetSocialMediaResponse, errs.Error) {
	var user dto.GetUsersSocialMedia
	var socialMedia dto.GetSocialMediaWithUser
	var resp dto.GetSocialMediaResponse

	rows, err := socialMediaPG.db.Query(getSocialMedias, userId)
	if err != nil {
		return nil, errs.NewInternalServerError(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&socialMedia.ID, &socialMedia.Name, &socialMedia.SocialMediaURL,
			&socialMedia.UserID, &socialMedia.CreatedAt, &socialMedia.UpdatedAt,
			&user.ID, &user.Username,
		)
		if err != nil {
			return nil, errs.NewInternalServerError(err.Error())
		}

		socialMedia.User = user
		resp.Data = append(resp.Data, socialMedia)
	}

	return &resp, nil
}

func (socialMediaPG *socialMediaPG) UpdateSocialMedia(sm dto.UpdateSocialMediaRequest) (*dto.UpdateSocialMediaResponse, errs.Error) {
	// name,social_media_url,id //request
	var resp dto.UpdateSocialMediaResponse

	err := socialMediaPG.db.QueryRow(updateSocialMedia, sm.Name, sm.SocialMediaURL, sm.ID).Scan(
		&resp.ID, &resp.Name, &resp.SocialMediaURL, &resp.UserID, &resp.UpdatedAt,
	)

	if err != nil {
		return nil, errs.NewInternalServerError(err.Error())
	}

	return &resp, nil
}

func (socialMediaPG *socialMediaPG) DeleteSocialMedia(id int) errs.Error {
	_, err := socialMediaPG.db.Exec(deleteSocialMedia, id)

	if err != nil {
		return errs.NewInternalServerError(err.Error())
	}

	return nil
}

func (socialMediaPG *socialMediaPG) GetUserID(id int) (int, errs.Error) {
	var userId int

	err := socialMediaPG.db.QueryRow(getUserId, id).Scan(
		&userId,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errs.NewBadRequest("Social Media Not Exist")
		}
		return 0, errs.NewInternalServerError(err.Error())
	}

	return userId, nil
}
