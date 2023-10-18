package service

import (
	"finalProject2/dto"
	"finalProject2/pkg/errs"
	"finalProject2/pkg/helpers"
	"finalProject2/repository/photo_repository"
)

type photoService struct {
	photoRepo photo_repository.Repository
}

type PhotoService interface {
	CreatePhoto(p dto.NewPhotoRequest) (*dto.NewPhotoResponse, errs.Error)
	GetPhotos() (*dto.GetPhotosResponse, errs.Error)
	UpdatePhotos(p dto.UpdatePhotoRequest) (*dto.UpdatePhotoResponse, errs.Error)
	DeletePhoto(p dto.DeletePhotoRequest) errs.Error
}

func NewPhotoService(photoRepo photo_repository.Repository) PhotoService {
	return &photoService{photoRepo: photoRepo}
}

func (ps *photoService) CreatePhoto(p dto.NewPhotoRequest) (*dto.NewPhotoResponse, errs.Error) {
	validateErr := helpers.ValidateStruct(&p)
	if validateErr != nil {
		return nil, validateErr
	}

	resp, err := ps.photoRepo.CreatePhoto(p)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ps *photoService) GetPhotos() (*dto.GetPhotosResponse, errs.Error) {
	resp, err := ps.photoRepo.GetPhotos()
	if err != nil {
		return nil, err
	}

	return resp, nil

}

func (ps *photoService) UpdatePhotos(p dto.UpdatePhotoRequest) (*dto.UpdatePhotoResponse, errs.Error) {
	belongsToAuthUser, err := ps.photoRepo.CheckIfPhotoBelongToUser(p.ID, p.UserID)
	if err != nil {
		return nil, err
	}

	if !belongsToAuthUser {
		return nil, errs.NewUnauthorizedError("You don't have permission to access this photo")
	}

	resp, err := ps.photoRepo.EditPhoto(p)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ps *photoService) DeletePhoto(p dto.DeletePhotoRequest) errs.Error {
	belongsToAuthUser, err := ps.photoRepo.CheckIfPhotoBelongToUser(p.ID, p.UserID)
	if err != nil {
		return err
	}

	if !belongsToAuthUser {
		return errs.NewUnauthorizedError("You don't have permission to access this photo")
	}

	err = ps.photoRepo.DeletePhoto(p.ID)
	if err != nil {
		return err
	}

	return nil
}
