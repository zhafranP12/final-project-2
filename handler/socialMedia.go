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

type SocialMediaHandler struct {
	socialMediaService service.SocialMediaService
}

func NewSocialMediaHandler(socialMediaService service.SocialMediaService) SocialMediaHandler {
	return SocialMediaHandler{socialMediaService: socialMediaService}
}

func (sh *SocialMediaHandler) CreateSocialMedia(c *gin.Context) {
	var newSocialMediaRequest dto.NewSocialMediaRequest

	if err := c.ShouldBindJSON(&newSocialMediaRequest); err != nil {
		errBindJson := errs.NewUnprocessibleEntityError("invalid json request body")
		c.AbortWithStatusJSON(errBindJson.Status(), errBindJson)
		return
	}

	jwtClaims := c.MustGet("user").(jwt.MapClaims)
	newSocialMediaRequest.UserID = int(jwtClaims["id"].(float64))
	resp, err := sh.socialMediaService.CreateSocialMedia(newSocialMediaRequest)

	if err != nil {
		c.AbortWithStatusJSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (sh *SocialMediaHandler) GetSocialMedias(c *gin.Context) {
	jwtClaims := c.MustGet("user").(jwt.MapClaims)
	userId := int(jwtClaims["id"].(float64))

	resp, err := sh.socialMediaService.GetSocialMedias(userId)
	if err != nil {
		c.AbortWithStatusJSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusCreated, resp)

}

func (sh *SocialMediaHandler) UpdateSocialMedia(c *gin.Context) {
	var updateSocialMediaReq dto.UpdateSocialMediaRequest

	if err := c.ShouldBindJSON(&updateSocialMediaReq); err != nil {
		errBindJson := errs.NewUnprocessibleEntityError("invalid json request body")
		c.AbortWithStatusJSON(errBindJson.Status(), errBindJson)
		return
	}

	jwtClaims := c.MustGet("user").(jwt.MapClaims)
	updateSocialMediaReq.UserID = int(jwtClaims["id"].(float64))

	idParam := c.Param("socialMediaId")
	idParamConvertToStr, errConv := strconv.Atoi(idParam)

	if errConv != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Error": errConv.Error(),
		})
		return
	}

	updateSocialMediaReq.ID = idParamConvertToStr
	resp, err := sh.socialMediaService.UpdateSocialMedia(updateSocialMediaReq)
	if err != nil {
		c.AbortWithStatusJSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (ph *SocialMediaHandler) DeleteSocialMedia(c *gin.Context) {
	jwtClaims := c.MustGet("user").(jwt.MapClaims)
	userId := int(jwtClaims["id"].(float64))

	idParam := c.Param("socialMediaId")
	idParamConvertToStr, errConv := strconv.Atoi(idParam)

	if errConv != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Error": errConv.Error(),
		})
		return
	}

	err := ph.socialMediaService.DeleteSocialMedia(idParamConvertToStr, userId)
	if err != nil {
		c.AbortWithStatusJSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, dto.DeleteSocialMediaResponse{
		Message: "Your Social Media Has Been Successfully Deleted",
	})

}
