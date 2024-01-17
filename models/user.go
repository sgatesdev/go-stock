package models

type User struct {
	ID       string `json:"id" gorm:"primary_key"`
	Email    string `json:"email" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
	// Stocks   []Stock `json:"stocks" gorm:"foreignKey:ID"`
	Username string `json:"username" gorm:"not null"`
}
