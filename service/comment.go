package service

import (
	"finalProject2/dto"
	"finalProject2/pkg/errs"
	"finalProject2/pkg/helpers"
	"finalProject2/repository/comment_repository"
)

type CommentService interface {
	CreateComment(c dto.NewCommentRequest) (*dto.NewCommentResponse, errs.Error)
	GetComments(userId int) (*dto.GetCommentsResponse, errs.Error)
	EditComment(c dto.UpdateCommentRequest) (*dto.UpdateCommentResponse, errs.Error)
	DeleteComment(c dto.DeleteCommentRequest) errs.Error
}

type commentService struct {
	commentRepo comment_repository.Repository
}

func NewCommentService(commentRepo comment_repository.Repository) CommentService {
	return &commentService{
		commentRepo: commentRepo,
	}
}

func (cs *commentService) CreateComment(c dto.NewCommentRequest) (*dto.NewCommentResponse, errs.Error) {
	validateErr := helpers.ValidateStruct(&c)
	if validateErr != nil {
		return nil, validateErr
	}

	photoExist, err := cs.commentRepo.PhotoExist(c.PhotoID)
	if err != nil {
		return nil, err
	}

	if !photoExist {
		return nil, errs.NewBadRequest("Photo Not Exist, You Can't Add Comment")
	}

	resp, err := cs.commentRepo.CreateComment(c)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (cs *commentService) GetComments(userId int) (*dto.GetCommentsResponse, errs.Error) {
	comments, err := cs.commentRepo.GetComments(userId)

	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (cs *commentService) EditComment(c dto.UpdateCommentRequest) (*dto.UpdateCommentResponse, errs.Error) {
	checkUserId, err := cs.commentRepo.GetUserId(c.ID)
	if err != nil {
		return nil, err
	}

	if checkUserId != c.UserID {
		return nil, errs.NewUnauthorizedError("You don't have permission to access this comment")
	}

	resp, err := cs.commentRepo.EditComment(c)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (cs *commentService) DeleteComment(c dto.DeleteCommentRequest) errs.Error {
	checkUserId, err := cs.commentRepo.GetUserId(c.ID)
	if err != nil {
		return err
	}

	if checkUserId != c.UserID {
		return errs.NewUnauthorizedError("You don't have permission to access this comment")
	}

	err = cs.commentRepo.DeleteComment(c.ID)
	if err != nil {
		return err
	}

	return nil

}
