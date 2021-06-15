package services

import (
	"errors"
	"log"

	"github.com/ardi7923/go-hello-api/models"
	"github.com/ardi7923/go-hello-api/repositories"
	"github.com/ardi7923/go-hello-api/requests"
	"github.com/mashingan/smapping"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	VerifyCredential(username string, password string) interface{}
	CreateUser(user requests.UserCreateRequest) models.User
	FindByUsername(username string) models.User
	IsDuplicateUsername(username string) bool
}

type authService struct {
	userRepository repositories.UserRepository
}

func NewAuthService(userRep repositories.UserRepository) AuthService {
	return &authService{
		userRepository: userRep,
	}
}

func (service *authService) VerifyCredential(username string, password string) interface{} {
	res := service.userRepository.VerifyCredential(username, password)

	if v, ok := res.(models.User); ok {
		comparedPassword := comparePassword(v.Password, []byte(password))
		if v.Username == username && comparedPassword {
			return res
		}
		return false
	}
	return false
}

func comparePassword(hashedPwd string, plainPassword []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (service *authService) CreateUser(user requests.UserCreateRequest) models.User {
	userToCreate := models.User{}
	err := smapping.FillStruct(&userToCreate, smapping.MapFields(&user))

	if err != nil {
		log.Fatal("Failed map", err)
	}
	res := service.userRepository.InsertUser(userToCreate)
	return res
}

func (service *authService) IsDuplicateUsername(username string) bool {
	res := service.userRepository.IsDuplicateUsername(username)
	return !errors.Is(res.Error, gorm.ErrRecordNotFound)
}

func (service *authService) FindByUsername(username string) models.User {
	return service.userRepository.FindByUsername(username)
}
