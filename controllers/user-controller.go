package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ardi7923/go-hello-api/helpers"
	"github.com/ardi7923/go-hello-api/requests"
	"github.com/ardi7923/go-hello-api/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	Update(context *gin.Context)
	Profile(context *gin.Context)
}

type userController struct {
	userService services.UserService
	jwtService  services.JWTService
}

func NewUserController(userService services.UserService, jwtService services.JWTService) UserController {
	return &userController{
		userService: userService,
		jwtService:  jwtService,
	}
}

// Update Profile godoc
// @Security ApiKeyAuth
// @Summary Update Profile User Actived
// @Desription Update Profile User Actived
// @ID Update Profile
// @Consume application/x-www-form-urlencoded
// @Produce json
// @Param id formData int true "User credentials"
// @Param name formData string true "User credentials"
// @Param username formData string true "User credentials"
// @Param password formData string true "User credentials"
// @Success 200 {object} helpers.Response
// @Failure 401 {object} helpers.Response
// @Router /user/profile [put]
// @Tags Profile
func (c *userController) Update(context *gin.Context) {
	var userUpdateRequest requests.UserUpdateRequest

	errRequest := context.ShouldBind(&userUpdateRequest)

	if (errRequest) != nil {
		res := helpers.BuildErrorResponse("Failed to proccess request", errRequest.Error(), helpers.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	id, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		panic(err.Error())
	}

	userUpdateRequest.Id = id

	proccess := c.userService.Update(userUpdateRequest)
	res := helpers.BuildResponse(true, "User has been Updated", proccess)
	context.JSON(http.StatusOK, res)
}

// Profile godoc
// @Security ApiKeyAuth
// @Summary Showing Data User Actived
// @Desription Showing Data User Actived
// @ID Profile
// @Consume json
// @Produce json
// @Success 200 {object} helpers.Response
// @Failure 401 {object} helpers.Response
// @Router /user/profile [get]
// @Tags Profile
func (c *userController) Profile(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	user := c.userService.Profile(fmt.Sprintf("%v", claims["user_id"]))
	res := helpers.BuildResponse(true, "OK", user)

	context.JSON(http.StatusOK, res)

}
