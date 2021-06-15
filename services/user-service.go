package services

import (
	"log"

	"github.com/ardi7923/go-hello-api/models"
	"github.com/ardi7923/go-hello-api/repositories"
	"github.com/ardi7923/go-hello-api/requests"
	"github.com/mashingan/smapping"
)

type UserService interface {
	Update(user requests.UserUpdateRequest) models.User
	Profile(userID string) models.User
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (service *userService) Update(user requests.UserUpdateRequest) models.User {
	userToUpdate := models.User{}
	err := smapping.FillStruct(&userToUpdate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}
	updatedUser := service.userRepository.UpdateUser(userToUpdate)
	return updatedUser
}

func (service *userService) Profile(userID string) models.User {
	return service.userRepository.ProfileUser(userID)
}
