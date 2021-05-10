package main

import (
	"github.com/gorilla/mux"
	"linebot/config"
	"linebot/controllers"
	"linebot/database"
	_ "linebot/database"
	"linebot/models"
	"net/http"
)

func main() {
	database.Db.AutoMigrate(&models.Food{}, &models.User{})
	router := mux.NewRouter()
	router.HandleFunc("/linebot", controllers.LineHandler)
	http.ListenAndServe(":"+config.Config.ServerPort, router)
}
