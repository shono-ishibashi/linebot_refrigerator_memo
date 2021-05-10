package models

import (
	"github.com/jinzhu/gorm"
	"linebot/database"
)

type User struct {
	gorm.Model
	UserId string
}

func FindUsersByUserId(user *User, userId string) {
	database.Db.Where("user_id = ?", userId).Find(&user)
}
