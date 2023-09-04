package models

import "time"

// change in num shares
type Event struct {
	ID           string `json:"id" gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	StockID      string `json:"stock_id" gorm:"foreignKey:StockID"`
	Stock        Stock
	Type         string `json:"type" gorm:"not null"`
	Date         string `json:"date" gorm:"not null"`
	SharesBefore int    `json:"shares_before" gorm:"not null"`
	SharesAfter  int    `json:"shares_after" gorm:"not null"`
}
