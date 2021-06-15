package controllers

import (
	"net/http"
	"strconv"

	"github.com/ardi7923/go-hello-api/helpers"
	"github.com/ardi7923/go-hello-api/models"
	"github.com/ardi7923/go-hello-api/requests"
	"github.com/ardi7923/go-hello-api/services"
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type authController struct {
	authService services.AuthService
	jwtService  services.JWTService
}

func NewAuthController(authService services.AuthService, jwtService services.JWTService) AuthController {
	return &authController{
		authService: authService,
		jwtService:  jwtService,
	}
}

// Login godoc
// @Summary Provides a JSON Web Token
// @Desription Authenticates a user and provider a JWT Authorize API calss
// @ID Login
// @Consume application/x-www-form-urlencoded
// @Produce json
// @Param username formData string true "User credentials"
// @Param password formData string true "User credentials"
// @Success 200 {object} helpers.Response
// @Failure 401 {object} helpers.Response
// @Router /auth/login [post]
// @Tags Authorization
func (c *authController) Login(ctx *gin.Context) {
	var loginRequest requests.LoginRequest
	errRequest := ctx.ShouldBind(&loginRequest)
	if errRequest != nil {
		res := helpers.BuildErrorResponse("Failed to proccess request", errRequest.Error(), helpers.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	authResult := c.authService.VerifyCredential(loginRequest.Username, loginRequest.Password)
	if v, ok := authResult.(models.User); ok {
		generateToken := c.jwtService.GenerateToken(strconv.FormatUint(v.Id, 10))
		v.Token = generateToken
		res := helpers.BuildResponse(true, "OK", v)
		ctx.JSON(http.StatusOK, res)
		return
	}
	res := helpers.BuildErrorResponse("Please check your credential", "Invalid Credential", helpers.EmptyObj{})
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)

}

// Register godoc
// @Summary Register User for Authenticates
// @Desription Authenticates a user and provider a JWT Authorize API calss
// @ID Register
// @Consume application/x-www-form-urlencoded
// @Produce json
// @Param name formData string true "User credentials"
// @Param username formData string true "User credentials"
// @Param password formData string true "User credentials"
// @Success 200 {object} helpers.Response
// @Failure 401 {object} helpers.Response
// @Router /auth/register [post]
// @Tags Authorization
func (c *authController) Register(ctx *gin.Context) {
	var registerRequest requests.UserCreateRequest
	errRequest := ctx.ShouldBind(&registerRequest)
	if errRequest != nil {
		res := helpers.BuildErrorResponse("Failed to process request", errRequest.Error(), helpers.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusConflict, res)
		return
	}

	if c.authService.IsDuplicateUsername(registerRequest.Username) {
		res := helpers.BuildErrorResponse("Failed to proccess request", "Duplicate Username", helpers.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
	} else {
		createdUser := c.authService.CreateUser(registerRequest)
		token := c.jwtService.GenerateToken(strconv.FormatUint(createdUser.Id, 10))
		createdUser.Token = token
		res := helpers.BuildResponse(true, "OK", createdUser)
		ctx.JSON(http.StatusCreated, res)
	}
}
