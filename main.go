package main

import (
	"linebot/config"
	"linebot/controllers"
	"linebot/database"
	_ "linebot/database"
	"linebot/models"
	"net/http"
)

func main() {
	database.Db.AutoMigrate(&models.Food{}, &models.User{})
	http.HandleFunc("/linebot", controllers.LineHandler)
	http.HandleFunc("/linebot/scheduler", controllers.SendMessageHandler)
	http.ListenAndServe(":"+config.Config.ServerPort, nil)
}
