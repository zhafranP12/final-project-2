package handler

import (
	"finalProject2/dto"
	"finalProject2/pkg/errs"
	"finalProject2/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type commentHandler struct {
	CommentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) commentHandler {
	return commentHandler{CommentService: commentService}
}

func (ch *commentHandler) CreateComment(c *gin.Context) {
	var comment dto.NewCommentRequest

	if err := c.ShouldBindJSON(&comment); err != nil {
		errBindJson := errs.NewUnprocessibleEntityError("invalid json request body")
		c.AbortWithStatusJSON(errBindJson.Status(), errBindJson)
		return
	}

	jwtClaims := c.MustGet("user").(jwt.MapClaims)
	comment.UserID = int(jwtClaims["id"].(float64))
	resp, err := ch.CommentService.CreateComment(comment)
	if err != nil {
		c.AbortWithStatusJSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (ch *commentHandler) GetComments(c *gin.Context) {

	jwtClaims := c.MustGet("user").(jwt.MapClaims)
	userId := int(jwtClaims["id"].(float64))

	resp, err := ch.CommentService.GetComments(userId)
	if err != nil {
		c.AbortWithStatusJSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (ch *commentHandler) UpdateComment(c *gin.Context) {
	var comment dto.UpdateCommentRequest

	if err := c.ShouldBindJSON(&comment); err != nil {
		errBindJson := errs.NewUnprocessibleEntityError("invalid json request body")
		c.AbortWithStatusJSON(errBindJson.Status(), errBindJson)
		return
	}

	jwtClaims := c.MustGet("user").(jwt.MapClaims)
	comment.UserID = int(jwtClaims["id"].(float64))

	idParam := c.Param("commentId")
	idParamConvertToStr, errConv := strconv.Atoi(idParam)
	comment.ID = idParamConvertToStr

	if errConv != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Error": errConv.Error(),
		})
		return
	}

	resp, err := ch.CommentService.EditComment(comment)
	if err != nil {
		c.AbortWithStatusJSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, resp)

}

func (commentHandler *commentHandler) DeleteComment(c *gin.Context) {
	var comment dto.DeleteCommentRequest

	jwtClaims := c.MustGet("user").(jwt.MapClaims)
	comment.UserID = int(jwtClaims["id"].(float64))

	idParam := c.Param("commentId")
	idParamConvertToStr, errConv := strconv.Atoi(idParam)

	if errConv != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Error": errConv.Error(),
		})
		return
	}

	comment.ID = idParamConvertToStr

	err := commentHandler.CommentService.DeleteComment(comment)
	if err != nil {
		c.AbortWithStatusJSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, dto.DeleteCommentResponse{
		Message: "Your Comment Has Been Succcessfully Deleted",
	})
}
