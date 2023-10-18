package photo_pg

import (
	"database/sql"
	"finalProject2/dto"
	"finalProject2/pkg/errs"
	"finalProject2/repository/photo_repository"
)

const (
	createPhoto = `
		INSERT INTO photos (user_id, title, caption, photo_url) 
		VALUES ($1, $2, $3, $4)
		RETURNING id,title,caption,photo_url,user_id,created_at
	`

	getPhotos = `
		SELECT photos.id, photos.title, photos.caption, 
		photos.photo_url,photos.user_id, 
		photos.created_at,photos.updated_at,
		users.email, users.username
		FROM photos LEFT JOIN users 
		ON users.id = photos.user_id
	`

	photoIdEqualToUserId = `
		select exists 
		(select true from photos where id = $1 AND user_id = $2)
	`

	getPhotoById = `
		SELECT user_id FROM photos WHERE id = $1
	`

	editPhoto = `
		UPDATE photos SET title = $1, caption = $2, photo_url = $3, updated_at = current_timestamp
		WHERE id = $4
		RETURNING id,title,caption,photo_url,user_id,updated_at
	`

	deletePhoto = `
		DELETE FROM photos WHERE id = $1
	`
)

type photoPG struct {
	db *sql.DB
}

func NewPhotoPG(db *sql.DB) photo_repository.Repository {
	return &photoPG{
		db: db,
	}
}

func (photoPG *photoPG) CreatePhoto(p dto.NewPhotoRequest) (*dto.NewPhotoResponse, errs.Error) {
	// req : user_id, title, caption, photo_url, user_id
	// res : id,title,caption,photo_url,user_id,created_at
	var resp dto.NewPhotoResponse
	err := photoPG.db.QueryRow(createPhoto, p.UserID, p.Title, p.Caption, p.PhotoURL).Scan(
		&resp.ID, &resp.Title, &resp.Caption, &resp.PhotoURL, &resp.UserID, &resp.CreatedAt,
	)

	if err != nil {
		return nil, errs.NewInternalServerError(err.Error())
	}

	return &resp, nil
}

func (photoPG *photoPG) GetPhotos() (*dto.GetPhotosResponse, errs.Error) {
	//  photos.id, photos.title, photos.caption,
	// 	photos.photo_url,photos.user_id,
	// 	photos.created_at,photos.updated_at,
	// 	users.email, users.username

	rows, err := photoPG.db.Query(getPhotos)
	if err != nil {
		return nil, errs.NewInternalServerError(err.Error())
	}

	defer rows.Close()

	var photo dto.GetPhotosWithUser
	var getPhotos dto.GetPhotosResponse

	for rows.Next() {
		err :=
			rows.Scan(
				&photo.ID, &photo.Title, &photo.Caption,
				&photo.PhotoURL, &photo.UserID, &photo.CreatedAt,
				&photo.UpdatedAt, &photo.User.Email, &photo.User.Username,
			)

		if err != nil {
			return nil, errs.NewInternalServerError(err.Error())
		}

		getPhotos.Data = append(getPhotos.Data, photo)
	}

	return &getPhotos, nil

}

func (photoPG *photoPG) EditPhoto(p dto.UpdatePhotoRequest) (*dto.UpdatePhotoResponse, errs.Error) {
	// title, caption ,photo_url ,id
	// id,title,caption,photo_url,user_id,updated_at
	var resp dto.UpdatePhotoResponse

	err := photoPG.db.QueryRow(editPhoto, &p.Title, &p.Caption, &p.PhotoURL, &p.ID).Scan(
		&resp.ID, &resp.Title, &resp.Caption, &resp.PhotoURL, &resp.UserID, &resp.UpdatedAt,
	)

	if err != nil {
		return nil, errs.NewInternalServerError(err.Error())
	}

	return &resp, nil
}

func (photoPG *photoPG) DeletePhoto(id int) errs.Error {
	_, err := photoPG.db.Exec(deletePhoto, id)
	if err != nil {
		return errs.NewInternalServerError(err.Error())
	}

	return nil
}

func (photoPG *photoPG) CheckIfPhotoBelongToUser(photoId, userId int) (bool, errs.Error) {
	var scanUserId int

	err := photoPG.db.QueryRow(getPhotoById, photoId).Scan(&scanUserId)
	if err != nil {

		return false, errs.NewInternalServerError(err.Error())
	}

	if scanUserId == userId {
		return true, nil
	}

	return false, nil
}
