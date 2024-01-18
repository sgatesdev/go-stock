package models

import "time"

// price of stock at a given time
type Holding struct {
	ID        string `json:"id" gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Shares    float32 `json:"shares" gorm:"column:shares;not null"`
	StockID   string  `json:"stock_id" gorm:"foreignKey:StockID"`
	Stock     Stock
	UserID    string `json:"user_id" gorm:"foreignKey:UserID"`
	User      User
}
