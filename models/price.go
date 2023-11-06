package models

import "time"

// price of stock at a given time
type Price struct {
	ID        string `json:"id" gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Received  int     `json:"t" gorm:"not null"`
	Type      string  `json:"type" gorm:"not null"`
	Price     float32 `json:"c" gorm:"column:price;not null"`
	StockID   string  `json:"stockId" gorm:"foreignKey:StockID"`
	Stock     Stock
}

// type can be: high, low, open, close

// from finnhub docs
// type Quote struct {
// 	// symbol
// 	S string `json:"symbol,omitempty"`
// 	// Open price of the day
// 	O float32 `json:"o,omitempty"`
// 	// High price of the day
// 	H float32 `json:"h,omitempty"`
// 	// Low price of the day
// 	L float32 `json:"l,omitempty"`
// 	// Current price
// 	C float32 `json:"c,omitempty"`
// 	// Previous close price
// 	Pc float32 `json:"pc,omitempty"`
// 	// timestamp
// 	T int `json:"t,omitempty"`
// }
