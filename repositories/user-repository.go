package repositories

import (
	"log"

	"github.com/ardi7923/go-hello-api/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	InsertUser(user models.User) models.User
	UpdateUser(user models.User) models.User
	VerifyCredential(username string, password string) interface{}
	IsDuplicateUsername(username string) (tx *gorm.DB)
	FindByUsername(username string) models.User
	ProfileUser(UserID string) models.User
}

type userConnection struct {
	connection *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{
		connection: db,
	}
}

func (db *userConnection) InsertUser(user models.User) models.User {
	user.Password = hashAndSalt([]byte(user.Password))
	db.connection.Save(&user)
	return user
}

func (db *userConnection) UpdateUser(user models.User) models.User {
	if user.Password != "" {
		user.Password = hashAndSalt([]byte(user.Password))
	} else {
		var tempUser models.User
		db.connection.Find(&tempUser, user.Id)
		user.Password = tempUser.Password
	}

	db.connection.Save(&user)
	return user
}

func (db *userConnection) VerifyCredential(username string, password string) interface{} {
	var user models.User
	res := db.connection.Where("username =  ?", username).Take(&user)
	if res.Error == nil {
		return user
	}

	return nil
}

func (db *userConnection) IsDuplicateUsername(username string) (tx *gorm.DB) {
	var user models.User
	return db.connection.Where("username = ?", username).Take(&user)
}

func (db *userConnection) FindByUsername(username string) models.User {
	var user models.User
	db.connection.Where("username = ?", username).Take(&user)
	return user
}

func (db *userConnection) ProfileUser(UserID string) models.User {
	var user models.User
	db.connection.Find(&user, UserID)
	return user
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		panic("Failed to hash a password")
	}

	return string(hash)
}
