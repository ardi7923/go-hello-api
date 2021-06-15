package models

type Project struct {
	Id          uint64 `gorm:"primary_key:auto_increment" json:"id"`
	Name        string `gorm:"type:varchar(191)" json:"name"`
	Duration    int8   `json:"duration"`
	Description string `gorm:"type:text" json:"description"`
}
