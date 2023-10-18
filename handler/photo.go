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

type photoHandler struct {
	PhotoService service.PhotoService
}

func NewPhotoHandler(photoService service.PhotoService) photoHandler {
	return photoHandler{
		PhotoService: photoService,
	}
}

func (ph *photoHandler) CreatePhoto(c *gin.Context) {
	var newPhotoRequest dto.NewPhotoRequest

	if err := c.ShouldBindJSON(&newPhotoRequest); err != nil {
		errBindJson := errs.NewUnprocessibleEntityError("invalid json request body")
		c.AbortWithStatusJSON(errBindJson.Status(), errBindJson)
		return
	}

	jwtClaims := c.MustGet("user").(jwt.MapClaims)
	newPhotoRequest.UserID = int(jwtClaims["id"].(float64))
	resp, err := ph.PhotoService.CreatePhoto(newPhotoRequest)

	if err != nil {
		c.AbortWithStatusJSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusCreated, resp)

}

func (ph *photoHandler) GetPhotos(c *gin.Context) {

	resp, err := ph.PhotoService.GetPhotos()
	if err != nil {
		c.AbortWithStatusJSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (ph *photoHandler) EditPhoto(c *gin.Context) {
	var updatePhotoRequest dto.UpdatePhotoRequest

	if err := c.ShouldBindJSON(&updatePhotoRequest); err != nil {
		errBindJson := errs.NewUnprocessibleEntityError("invalid json request body")
		c.AbortWithStatusJSON(errBindJson.Status(), errBindJson)
		return
	}

	jwtClaims := c.MustGet("user").(jwt.MapClaims)
	updatePhotoRequest.UserID = int(jwtClaims["id"].(float64))

	idParam := c.Param("photoId")
	idParamConvertToStr, errConv := strconv.Atoi(idParam)

	if errConv != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Error": errConv.Error(),
		})
		return
	}

	updatePhotoRequest.ID = idParamConvertToStr
	resp, err := ph.PhotoService.UpdatePhotos(updatePhotoRequest)
	if err != nil {
		c.AbortWithStatusJSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, resp)

}

func (ph *photoHandler) DeletePhoto(c *gin.Context) {
	var deletePhotoRequest dto.DeletePhotoRequest

	jwtClaims := c.MustGet("user").(jwt.MapClaims)
	deletePhotoRequest.UserID = int(jwtClaims["id"].(float64))

	idParam := c.Param("photoId")
	idParamConvertToStr, errConv := strconv.Atoi(idParam)

	if errConv != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Error": errConv.Error(),
		})
		return
	}

	deletePhotoRequest.ID = idParamConvertToStr
	err := ph.PhotoService.DeletePhoto(deletePhotoRequest)
	if err != nil {
		c.AbortWithStatusJSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, dto.DeletePhotoResponse{
		Message: "Your Photo Has Been Successfully Deleted",
	})
}
