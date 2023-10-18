package handler

import (
	"finalProject2/dto"
	"finalProject2/pkg/errs"
	"finalProject2/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type userHandler struct {
	UserService service.UserService
}

func NewUserHandler(userService service.UserService) userHandler {
	return userHandler{
		UserService: userService,
	}
}

func (us *userHandler) register(c *gin.Context) {
	var newUserRequest dto.NewUserRequest

	if err := c.ShouldBindJSON(&newUserRequest); err != nil {
		errBindJson := errs.NewUnprocessibleEntityError("invalid json request body")
		c.AbortWithStatusJSON(errBindJson.Status(), errBindJson)
		return
	}

	resp, err := us.UserService.CreateUser(newUserRequest)
	if err != nil {
		c.AbortWithStatusJSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (us *userHandler) Login(c *gin.Context) {
	var user dto.LoginRequest

	if err := c.ShouldBindJSON(&user); err != nil {
		errBindJson := errs.NewUnprocessibleEntityError("invalid json request body")
		c.AbortWithStatusJSON(errBindJson.Status(), errBindJson)
		return
	}

	resp, err := us.UserService.Login(user)
	if err != nil {
		c.AbortWithStatusJSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (us *userHandler) EditUser(c *gin.Context) {
	var user dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		errBindJson := errs.NewUnprocessibleEntityError("invalid json request body")
		c.AbortWithStatusJSON(errBindJson.Status(), errBindJson)
		return
	}

	jwtClaims := c.MustGet("user").(jwt.MapClaims)
	user.ID = int(jwtClaims["id"].(float64))

	resp, err := us.UserService.UpdateUser(user)
	if err != nil {
		c.AbortWithStatusJSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, resp)

}

func (us *userHandler) DeleteUser(c *gin.Context) {

	jwtClaims := c.MustGet("user").(jwt.MapClaims)
	id := int(jwtClaims["id"].(float64))

	err := us.UserService.DeleteUser(id)
	if err != nil {
		c.AbortWithStatusJSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, dto.DeleteUserResponse{Message: "Your Account Has Been Successfully Deleted"})

}
