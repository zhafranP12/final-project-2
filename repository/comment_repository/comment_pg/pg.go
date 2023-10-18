package comment_pg

import (
	"database/sql"
	"finalProject2/dto"
	"finalProject2/pkg/errs"
	"finalProject2/repository/comment_repository"
	"fmt"
)

type commentPG struct {
	db *sql.DB
}

func NewCommentPG(db *sql.DB) comment_repository.Repository {
	return &commentPG{
		db: db,
	}
}

const (
	createComment = `
		INSERT INTO comments (message,photo_id,user_id) 
		VALUES ($1,$2,$3)
		RETURNING id,message,photo_id,user_id,created_at
	`

	photoExist = `
		SELECT COUNT(1) FROM photos WHERE id = $1 
	`

	getComments = `
		SELECT 
			comments.id,comments.message,comments.photo_id,
			comments.user_id,comments.updated_at,
			comments.created_at,
			photos.id,photos.title,photos.caption,
			photos.photo_url,
			users.email,users.username
		FROM comments 
		LEFT JOIN photos 
			ON comments.photo_id = photos.id
		LEFT JOIN users
			ON photos.user_id = users.id
		WHERE comments.user_id = $1
	`

	updateComment = `
		UPDATE comments SET message = $1, updated_at = current_timestamp WHERE id = $2
	`

	getUserIdFromPhoto = `
		SELECT user_id FROM photos WHERE id = $1	
	`

	getUserIdFromComment = `
		SELECT user_id FROM comments WHERE id = $1
	`

	updateResponse = `
	SELECT 
		comments.id, photos.title, photos.caption,
		photos.photo_url,photos.user_id,comments.updated_at
	FROM comments LEFT JOIN photos
		ON comments.photo_id = photos.id 
		WHERE comments.id = $1
	`

	deleteComment = `
		DELETE FROM comments WHERE id = $1
	`
)

func (commentPG *commentPG) CreateComment(c dto.NewCommentRequest) (*dto.NewCommentResponse, errs.Error) {
	// message,photo_id,user_id
	// id,message,photo_id,user_id,created_at

	var resp dto.NewCommentResponse

	err := commentPG.db.QueryRow(createComment, c.Message, c.PhotoID, c.UserID).Scan(
		&resp.ID, &resp.Message, &resp.PhotoID, &resp.UserID, &resp.CreatedAt,
	)

	if err != nil {
		return nil, errs.NewInternalServerError(err.Error())
	}

	return &resp, nil

}

func (commentPG *commentPG) GetUserId(id int) (int, errs.Error) {
	// if comment not found, return 404 error
	var userId int

	err := commentPG.db.QueryRow(getUserIdFromComment, id).Scan(
		&userId,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errs.NewBadRequest("Comment Not Exist")
		}
		return 0, errs.NewInternalServerError(err.Error())
	}

	return userId, nil
}

func (commentPG *commentPG) PhotoExist(photoId int) (bool, errs.Error) {
	var count int

	err := commentPG.db.QueryRow(photoExist, photoId).Scan(
		&count,
	)

	if err != nil {
		return false, errs.NewInternalServerError(err.Error())
	}

	return count > 0, nil
}

func (commentPG *commentPG) GetComments(userId int) (*dto.GetCommentsResponse, errs.Error) {

	var user dto.GetUsersComment
	var photo dto.GetPhotosComment
	var comment dto.GetCommentsWithUserAndPhoto
	var resp dto.GetCommentsResponse

	rows, err := commentPG.db.Query(getComments, userId)

	if err != nil {
		return nil, errs.NewInternalServerError(err.Error())
	}

	defer rows.Close()

	photo.UserID = userId
	user.ID = userId

	for rows.Next() {
		err := rows.Scan(
			&comment.ID, &comment.Message, &comment.PhotoID,
			&comment.UserID, &comment.UpdatedAt, &comment.CreatedAt,
			&photo.ID, &photo.Title, &photo.Caption, &photo.PhotoURL,
			&user.Email, &user.Username,
		)

		if err != nil {
			return nil, errs.NewInternalServerError(err.Error())
		}

		comment.Photo = photo
		comment.User = user

		fmt.Println(comment)

		resp.Data = append(resp.Data, comment)

	}

	return &resp, nil

}

func (commentPG *commentPG) EditComment(c dto.UpdateCommentRequest) (*dto.UpdateCommentResponse, errs.Error) {
	var resp dto.UpdateCommentResponse

	_, err := commentPG.db.Exec(updateComment, c.Message, c.ID)
	if err != nil {
		return nil, errs.NewInternalServerError(err.Error())
	}

	fmt.Println("test", c.ID)
	err = commentPG.db.QueryRow(updateResponse, c.ID).Scan(
		&resp.ID, &resp.Title, &resp.Caption,
		&resp.PhotoURL, &resp.UserID, &resp.UpdatedAt,
	)
	if err != nil {
		return nil, errs.NewInternalServerError(err.Error())
	}

	return &resp, nil
}

func (commentPG *commentPG) DeleteComment(id int) errs.Error {

	_, err := commentPG.db.Exec(deleteComment, id)
	if err != nil {
		return errs.NewInternalServerError(err.Error())
	}

	return nil

}
