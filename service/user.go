package service

import (
	"finalProject2/dto"
	"finalProject2/pkg/errs"
	"finalProject2/pkg/helpers"
	user_repository "finalProject2/repository"
)

type userService struct {
	UserRepo user_repository.Repository
}

type UserService interface {
	CreateUser(newUser dto.NewUserRequest) (*dto.NewUserResponse, errs.Error)
	Login(u dto.LoginRequest) (*dto.LoginResponse, errs.Error)
	UpdateUser(u dto.UpdateUserRequest) (*dto.UpdateUserResponse, errs.Error)
	DeleteUser(id int) errs.Error
}

func NewUserService(userRepo user_repository.Repository) UserService {
	return &userService{
		UserRepo: userRepo,
	}
}

func (userService *userService) CreateUser(newUser dto.NewUserRequest) (*dto.NewUserResponse, errs.Error) {
	validateErr := helpers.ValidateStruct(&newUser)
	if validateErr != nil {
		return nil, validateErr
	}

	countEmail, err := userService.UserRepo.CountEmail(newUser.Email)
	if err != nil {
		return nil, err
	}

	if countEmail > 0 {
		return nil, errs.NewBadRequest("Email Has Already Been Used")
	}

	countUsername, err := userService.UserRepo.CountUsername(newUser.Username)
	if err != nil {
		return nil, err
	}

	if countUsername > 0 {
		return nil, errs.NewBadRequest("Username Has Already Been Used")
	}

	generatePW, errPW := helpers.GenerateHashedPassword([]byte(newUser.Password))
	if errPW != nil {
		return nil, errs.NewInternalServerError(errPW.Error())
	}

	newUser.Password = generatePW
	user, errs := userService.UserRepo.CreateUser(newUser)

	if errs != nil {
		return nil, errs
	}

	return user, nil
}

func (userService *userService) Login(u dto.LoginRequest) (*dto.LoginResponse, errs.Error) {
	validateErr := helpers.ValidateStruct(&u)
	if validateErr != nil {
		return nil, validateErr
	}

	user, err := userService.UserRepo.Login(u.Email)
	if err != nil {
		return nil, err
	}

	var resp dto.LoginResponse
	compare := helpers.ComparePass([]byte(user.Password), []byte(u.Password))
	if compare {
		token, errService := helpers.GenerateToken(user.ID, user.Email)

		if errService != nil {
			return nil, errs.NewInternalServerError(errService.Error())
		}

		resp.Token = token
	} else {
		return nil, errs.NewUnauthenticatedError("Kombinasi Email dan Password Salah")
	}

	return &resp, nil

}

func (userService *userService) UpdateUser(u dto.UpdateUserRequest) (*dto.UpdateUserResponse, errs.Error) {
	validateErr := helpers.ValidateStruct(&u)
	if validateErr != nil {
		return nil, validateErr
	}

	countEmail, err := userService.UserRepo.CountEmail(u.Email)
	if err != nil {
		return nil, err
	}

	if countEmail > 0 {
		return nil, errs.NewBadRequest("Email Has Already Been Used")
	}

	countUsername, err := userService.UserRepo.CountUsername(u.Username)
	if err != nil {
		return nil, err
	}

	if countUsername > 0 {
		return nil, errs.NewBadRequest("Username Has Already Been Used")
	}

	user, err := userService.UserRepo.EditUser(u)
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (userService *userService) DeleteUser(id int) errs.Error {
	err := userService.UserRepo.DeleteUser(id)
	if err != nil {
		return err
	}

	return nil
}
