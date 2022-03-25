package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"linebot/config"
	"log"
)

var Db *gorm.DB

var err error

func init() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=Asia/Tokyo",
		config.Config.DbHost,
		config.Config.DbUserName,
		config.Config.DbUserPassword,
		config.Config.DbName,
		config.Config.DbPort,
	)

	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("Successfully connect database")
	}
}
