package models

type User struct {
	Id       uint64 `gorm:"primary_key:auto_increment" json:"id"`
	Name     string `gorm:"type:varchar(191)" json:"name"`
	Username string `gorm:"uniqueIndex;type:varchar(191)" json:"username"`
	Password string `gorm:"->;<-;not null" json:"-"`
	Token    string `gorm:"-" json:"token,omitempty"`
}
