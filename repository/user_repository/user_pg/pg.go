package user_pg

import (
	"database/sql"
	"errors"
	"finalProject2/dto"
	"finalProject2/entity"
	"finalProject2/pkg/errs"
	user_repository "finalProject2/repository"
)

const (
	createUser = `
		INSERT INTO users (username,email,age,password) 
		VALUES ($1, $2, $3, $4)
		RETURNING age,email,id,username
	`

	login = `
		SELECT id,password,email FROM users WHERE email = $1
	`

	checkIfUserExist = `
		SELECT COUNT(*) FROM users WHERE id = $1; 
	`

	editUser = `
		UPDATE users SET email = $1, username = $2, updated_at = current_timestamp  WHERE id = $3
		RETURNING id,email,username,age,updated_at
	`

	deleteUser = `
		DELETE FROM users WHERE id = $1
	`

	countEmail = `
		SELECT COUNT(1) FROM users WHERE email = $1
	`

	countUsername = `
		SELECT COUNT(1) FROM users WHERE username = $1
	`
)

type userPG struct {
	db *sql.DB
}

func NewOrderPG(db *sql.DB) user_repository.Repository {
	return &userPG{db: db}
}

func (userPG *userPG) CreateUser(newUser dto.NewUserRequest) (*dto.NewUserResponse, errs.Error) {
	// username,email,age,password
	var user dto.NewUserResponse
	err := userPG.db.QueryRow(createUser, newUser.Username, newUser.Email, newUser.Age, newUser.Password).Scan(
		&user.Age, &user.Email, &user.ID, &user.Username,
	)

	if err != nil {
		return nil, errs.NewInternalServerError(err.Error())
	}

	return &user, nil

}

func (userPG *userPG) Login(email string) (*entity.User, errs.Error) {
	var u entity.User

	err := userPG.db.QueryRow(login, email).Scan(&u.ID, &u.Password, &u.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("user not found")
		}
		return nil, errs.NewInternalServerError(err.Error())
	}

	return &u, nil
}

func (userPG *userPG) EditUser(user dto.UpdateUserRequest) (*dto.UpdateUserResponse, errs.Error) {
	// id,email,username,age,updated_at
	var u dto.UpdateUserResponse
	err := userPG.db.QueryRow(editUser, &user.Email, &user.Username, &user.ID).Scan(&u.ID, &u.Email, &u.Username, &u.Age, &u.UpdatedAt)
	if err != nil {
		return nil, errs.NewInternalServerError(err.Error())
	}

	return &u, nil
}

func (userPG *userPG) DeleteUser(id int) errs.Error {

	_, err := userPG.db.Exec(deleteUser, id)
	if err != nil {
		return errs.NewInternalServerError(err.Error())
	}

	return nil
}

func (userPG *userPG) CountEmail(email string) (int, errs.Error) {
	var count int

	err := userPG.db.QueryRow(countEmail, email).Scan(&count)
	if err != nil {
		return 0, errs.NewInternalServerError(err.Error())
	}

	return count, nil
}

func (userPG *userPG) CountUsername(username string) (int, errs.Error) {
	var count int

	err := userPG.db.QueryRow(countUsername, username).Scan(&count)
	if err != nil {
		return 0, errs.NewInternalServerError(err.Error())
	}

	return count, nil
}
