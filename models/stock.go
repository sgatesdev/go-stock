package models

import "time"

// tracks symbols across userbase to build finnhub queries
type Stock struct {
	ID        string `json:"id" gorm:"primarykey"`
	Name      string `json:"name" gorm:"column:name;not null"`
	Symbol    string `json:"symbol" gorm:"column:symbol;not null"`
	Poll      bool   `json:"poll" gorm:"column:poll;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Stocks struct {
	Stocks []Stock `json:"stocks"`
}
