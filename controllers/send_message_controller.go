package controllers

import (
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
	"linebot/line_utils"
	"linebot/models"
	"net/http"
)

func SendMessageHandler(w http.ResponseWriter, _ *http.Request) {
	var userIds []string
	models.FindUserIdByExpirationDate(&userIds)

	for _, userId := range userIds {
		var foods []models.Food
		models.FindFoodsByUserIdAndExpirationDate(&foods, userId)

		if len(foods) == 0 {
			return
		}
		message := "期限が近い食品があります"
		for _, food := range foods {
			message += "\n" + generateFoodMessageFormat(food)
		}

		fmt.Println(userIds)
		fmt.Println(message)

		line_utils.Bot.PushMessage(userId, linebot.NewTextMessage(message)).Do()
	}

	w.WriteHeader(http.StatusOK)
}

func generateFoodMessageFormat(food models.Food) string {
	stringDate := food.ExpirationDate.Format(line_utils.DateFormat)
	message := fmt.Sprintf("%s : %s", food.Name, stringDate)
	return message
}
