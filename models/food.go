package models

import (
	"gorm.io/gorm"
	"linebot/database"
	"time"
)

const (
	InStockStatus   = "0"
	AteStatus       = "1"
	DiscardedStatus = "2"
)

type Food struct {
	gorm.Model
	UserId         string    `gorm:"type:varchar(100);not null"`
	Name           string    `gorm:"type:varchar(100)"`
	ExpirationDate time.Time `gorm:"not null;type:date"`
	Status         string    `gorm:"type:varchar(1);not null"`
}

type FoodRate struct {
	AteStatusCount       int
	DiscardedStatusCount int
	AteRate              int
}

func (food Food) InsertFood() {
	database.Db.Create(&food)
}

func (food Food) UpdateFood() {
	database.Db.Save(&food)
}

func (food Food) DeleteFood() {
	database.Db.Delete(&food)
}

func FindFoodByUserIdAndStatus(foods *[]Food, UserId string, status string) {
	database.Db.Where("user_id = ? AND status = ?", UserId, status).Order("expiration_date").Find(foods)
}

func FindFoodByFoodId(food *Food, foodId uint) {
	database.Db.First(food, foodId)
}

func FindRate(foodRate *FoodRate, UserId string) {
	var query = "SELECT COUNT(status = ? OR null) AS ate_status_count, COUNT(status = ? OR null) AS discarded_status_count FROM foods WHERE user_id = ?"
	database.Db.Raw(query, AteStatus, DiscardedStatus, UserId).Scan(foodRate)
}

func FindUserIdByExpirationDate(UserIds *[]string) {
	//query := "SELECT DISTINCT user_id as user_ids FROM foods WHERE expiration_date - interval '2day' <=  current_date AND status = ? AND deleted_at IS NULL ORDER BY user_id"
	//database.Db.Raw(query, InStockStatus).Scan(UserIds)
	database.Db.Model(&Food{}).Distinct().Where("expiration_date - interval '1day' <=  current_date AND status = ?", InStockStatus).Pluck("UserId", UserIds)
}

func FindFoodsByUserIdAndExpirationDate(foods *[]Food, UserId string) {
	database.Db.Where("user_id = ? AND expiration_date - interval '1day' <=  current_date AND expiration_date >=  current_date AND status = ? ", UserId, InStockStatus).Find(foods)
}
