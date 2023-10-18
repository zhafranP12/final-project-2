package service

import (
	"finalProject2/dto"
	"finalProject2/pkg/errs"
	social_media_repository "finalProject2/repository/socialMedia_repository"
)

type SocialMediaService interface {
	CreateSocialMedia(s dto.NewSocialMediaRequest) (*dto.NewSocialMediaResponse, errs.Error)
	GetSocialMedias(userId int) (*dto.GetSocialMediaResponse, errs.Error)
	UpdateSocialMedia(s dto.UpdateSocialMediaRequest) (*dto.UpdateSocialMediaResponse, errs.Error)
	DeleteSocialMedia(id, userId int) errs.Error
}

func NewSocialMediaService(socialMediaRepo social_media_repository.Repository) SocialMediaService {
	return &socialMediaService{socialMediaRepo: socialMediaRepo}
}

type socialMediaService struct {
	socialMediaRepo social_media_repository.Repository
}

func (sc *socialMediaService) CreateSocialMedia(s dto.NewSocialMediaRequest) (*dto.NewSocialMediaResponse, errs.Error) {
	resp, err := sc.socialMediaRepo.CreateSocialMedia(s)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (sc *socialMediaService) GetSocialMedias(userId int) (*dto.GetSocialMediaResponse, errs.Error) {

	resp, err := sc.socialMediaRepo.GetSocialMedias(userId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (sc *socialMediaService) UpdateSocialMedia(s dto.UpdateSocialMediaRequest) (*dto.UpdateSocialMediaResponse, errs.Error) {

	checkUserId, err := sc.socialMediaRepo.GetUserID(s.ID)
	if err != nil {
		return nil, err
	}

	if checkUserId != s.UserID {
		return nil, errs.NewUnauthorizedError("You don't have permission to access this social media")
	}

	resp, err := sc.socialMediaRepo.UpdateSocialMedia(s)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (sc *socialMediaService) DeleteSocialMedia(id, userId int) errs.Error {
	checkUserId, err := sc.socialMediaRepo.GetUserID(id)
	if err != nil {
		return err
	}

	if checkUserId != userId {
		return errs.NewUnauthorizedError("You don't have permission to access this social media")
	}

	err = sc.socialMediaRepo.DeleteSocialMedia(id)
	if err != nil {
		return err
	}

	return nil
}
