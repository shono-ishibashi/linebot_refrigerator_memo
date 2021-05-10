package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"linebot/config"
	"linebot/utils"
	"log"
)

var Db *gorm.DB

var err error

func init() {
	dbConnectInfo := fmt.Sprintf(
		`postgres://%s:%s@%s:%s/%s`,
		config.Config.DbUserName,
		config.Config.DbUserPassword,
		config.Config.DbHost,
		config.Config.DbPort,
		config.Config.DbName,
	)

	fmt.Println(dbConnectInfo)

	Db, err = gorm.Open(config.Config.DbDriverName, dbConnectInfo)
	Db.SetLogger(utils.Logger)

	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println("Successfully connect database")
	}
}
