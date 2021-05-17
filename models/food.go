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

// UpdateFood 新規に保存する。
func (food Food) InsertFood() {
	database.Db.Create(&food)
}

// UpdateFood 更新する
func (food Food) UpdateFood() {
	database.Db.Save(&food)
}

// DeleteFood 削除する
func (food Food) DeleteFood() {
	database.Db.Delete(&food)
}

// FindFoodByUserIdAndStatus ユーザーIDとStatusでFoodを検索する。
func FindFoodByUserIdAndStatus(foods *[]Food, UserId string, status string) {
	database.Db.Where("user_id = ? AND status = ?", UserId, status).Order("expiration_date").Find(foods)
}

// FindFoodByFoodId 在庫状態のFoodをFood IDで検索する。
func FindFoodByFoodId(food *Food, foodId uint) {
	database.Db.Where("status = ?", InStockStatus).First(food, foodId)
}

// FindRate ユーザーの食べたカウント、破棄したカウントを取得する。
func FindRate(foodRate *FoodRate, UserId string) {
	var query = "SELECT COUNT(status = ? OR null) AS ate_status_count, COUNT(status = ? OR null) AS discarded_status_count FROM foods WHERE user_id = ?"
	database.Db.Raw(query, AteStatus, DiscardedStatus, UserId).Scan(foodRate)
}

// FindUserIdByExpirationDate 期限が2日以内の食品を持つUserIdを検索する。
func FindUserIdByExpirationDate(UserIds *[]string) {
	database.Db.Model(&Food{}).Distinct().Where("expiration_date - interval '1day' <=  current_date AND status = ?", InStockStatus).Pluck("UserId", UserIds)
}

// FindFoodsByUserIdAndExpirationDate 期限が2日以内の食品をUserIdで検索する。
func FindFoodsByUserIdAndExpirationDate(foods *[]Food, UserId string) {
	database.Db.Where("user_id = ? AND expiration_date - interval '1day' <=  current_date AND expiration_date >=  current_date AND status = ? ", UserId, InStockStatus).Find(foods)
}
