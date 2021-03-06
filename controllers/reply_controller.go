package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
	"linebot/line_utils"
	"linebot/models"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

func LineHandler(w http.ResponseWriter, r *http.Request) {
	events, err := line_utils.Bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				SentMessage := message.Text

				// reply food list in
				if SentMessage == "list" {
					replayFoodList(line_utils.Bot, event)
					return
				}

				if SentMessage == "rate" {
					replyFoodsEatenRate(line_utils.Bot, event)
					return
				}

				// reply add form
				if _, err = line_utils.Bot.ReplyMessage(event.ReplyToken,
					linebot.NewFlexMessage(SentMessage, line_utils.GenerateAddFoodConfirmationTemplate(SentMessage))).Do(); err != nil {
					log.Fatalln(err)
				}
			default:
				if _, err = line_utils.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("商品名を入力してください")).Do(); err != nil {
					log.Fatalln(err)
				}
			}

		case linebot.EventTypePostback:
			var param map[string]interface{}
			err = json.Unmarshal([]byte(event.Postback.Data), &param)

			if err != nil {
				log.Fatalln(err)
			}

			switch param["type"] {
			case "add":
				expirationDate, err := time.Parse(
					line_utils.DateFormat,
					event.Postback.Params.Date,
				)

				if err != nil {
					log.Fatalln(err)
				}

				food := models.Food{
					Name:           param["foodName"].(string),
					ExpirationDate: expirationDate,
					Status:         models.InStockStatus,
					UserId:         event.Source.UserID,
				}
				addFood(line_utils.Bot, event, &food)
			case "detail":
				foodId := convertStringToUint(param["foodId"].(string))
				replyFoodDetail(line_utils.Bot, event, foodId)
			case "eat":
				foodId := convertStringToUint(param["foodId"].(string))
				replyEatFood(line_utils.Bot, event, foodId)
			case "discarded":
				foodId := convertStringToUint(param["foodId"].(string))
				replyDiscardFood(line_utils.Bot, event, foodId)
			case "delete":
				foodId := convertStringToUint(param["foodId"].(string))
				replyDeleteFood(line_utils.Bot, event, foodId)
			}
		}
	}
}

func addFood(bot *linebot.Client, event *linebot.Event, food *models.Food) {
	food.InsertFood()
	replayMessage := fmt.Sprintf("%s をLINE冷蔵庫に追加しました。\n期限: %s", food.Name, food.ExpirationDate.Format(line_utils.DateFormat))
	_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replayMessage)).Do()

	if err != nil {
		log.Fatalln(err)
	}
}

func replayFoodList(bot *linebot.Client, event *linebot.Event) {
	var foods []models.Food
	models.FindFoodByUserIdAndStatus(&foods, event.Source.UserID, models.InStockStatus)

	if len(foods) == 0 {
		replyMessage := "冷蔵庫の中身が空だよ。食品名を送信して冷蔵庫に追加してね！"
		_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	replayFlex := line_utils.GenerateListTemplate(foods)
	_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewFlexMessage("冷蔵庫の中身一覧だよ", replayFlex)).Do()

	if err != nil {
		log.Fatalln(err)
	}
}

func replyFoodDetail(bot *linebot.Client, event *linebot.Event, foodId uint) {
	var food models.Food
	models.FindFoodByFoodId(&food, foodId)
	replayFlex := line_utils.GenerateDetailTemplate(food)
	_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewFlexMessage(food.Name, replayFlex)).Do()

	if err != nil {
		log.Fatalln(err)
	}
}

func replyEatFood(bot *linebot.Client, event *linebot.Event, foodId uint) {
	var food models.Food
	models.FindFoodByFoodId(&food, foodId)
	food.Status = models.AteStatus
	food.UpdateFood()
	replyMessage := fmt.Sprintf("「%s」をから冷蔵庫から食べました！", food.Name)
	_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()

	if err != nil {
		log.Fatalln(err)
	}
}

func replyDiscardFood(bot *linebot.Client, event *linebot.Event, foodId uint) {
	var food models.Food
	models.FindFoodByFoodId(&food, foodId)
	food.Status = models.DiscardedStatus
	food.UpdateFood()
	replyMessage := fmt.Sprintf("「%s」を破棄しました、、、。", food.Name)
	_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()

	if err != nil {
		log.Fatalln(err)
	}
}
func replyDeleteFood(bot *linebot.Client, event *linebot.Event, foodId uint) {
	var food models.Food
	models.FindFoodByFoodId(&food, foodId)
	food.DeleteFood()
	replyMessage := fmt.Sprintf("「%s」を削除しました。", food.Name)
	_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()

	if err != nil {
		log.Fatalln(err)
	}
}

func replyFoodsEatenRate(bot *linebot.Client, event *linebot.Event) {
	var foodRate models.FoodRate
	models.FindRate(&foodRate, event.Source.UserID)
	if foodRate.AteStatusCount == 0 && foodRate.DiscardedStatusCount == 0 {
		replyMessage := fmt.Sprintf("データがありません。")
		_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()
		if err != nil {
			log.Fatalln(err)
		}
		return
	}
	foodRate.AteRate = int(math.Round(float64(foodRate.AteStatusCount) * 100.0 / (float64(foodRate.AteStatusCount) + float64(foodRate.DiscardedStatusCount))))
	replayFlex := line_utils.GenerateFoodsEatenRateTemplate(foodRate)
	altText := "食べた割合"

	_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewFlexMessage(altText, replayFlex)).Do()

	if err != nil {
		log.Fatalln(err)
	}
}

func convertStringToUint(s string) uint {
	value, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		log.Fatalln(err)
	}
	return uint(value)
}
