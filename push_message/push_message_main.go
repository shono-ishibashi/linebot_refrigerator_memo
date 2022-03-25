package main

import (
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/robfig/cron/v3"
	"linebot/line_utils"
	"linebot/models"
	"runtime"
)

func main() {
	c := cron.New()
	_, err := c.AddFunc("0 8 * * *", SendMessage)

	if err != nil {
		fmt.Println(err)
	}

	c.Start()
	runtime.Goexit()
}

func SendMessage() {
	var userIds []string
	models.FindUserIds(&userIds)

	for _, userId := range userIds {
		sendExpirationDateMessage(userId)
		sendExpiredFoodMessage(userId)
	}

}

func sendExpirationDateMessage(userId string) {
	var foods []models.Food
	models.FindFoodsByUserIdAndExpirationDate(&foods, userId)

	if len(foods) == 0 {
		return
	}
	message := "期限が近い食品があります"
	for _, food := range foods {
		message += "\n" + generateExpirationDateMessageFormat(food)
	}

	line_utils.Bot.PushMessage(userId, linebot.NewTextMessage(message)).Do()
}

func sendExpiredFoodMessage(userId string) {
	var foods []models.Food
	models.FindExpiredFood(&foods, userId)

	if len(foods) == 0 {
		return
	}
	message := "期限切れの食品があります"
	for _, food := range foods {
		message += "\n" + generateExpirationDateMessageFormat(food)
	}

	line_utils.Bot.PushMessage(userId, linebot.NewTextMessage(message)).Do()

}

func generateExpirationDateMessageFormat(food models.Food) string {
	stringDate := food.ExpirationDate.Format(line_utils.DateFormat)
	message := fmt.Sprintf("%s : %s", food.Name, stringDate)
	return message
}
