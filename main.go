package main

import (
	"net/http"

	handlers "samgates.io/go-stock/handlers"
	models "samgates.io/go-stock/models"
	scheduler "samgates.io/go-stock/scheduler"
	"samgates.io/go-stock/utils"
)

func main() {
	// connect to db
	db := utils.SetupDB()

	stockHandler := handlers.NewStockHandler(db)
	stockHandler.RegisterRoutes()

	// fetch stocks
	stocks := []models.Stock{}
	err := db.Find(&stocks).Where("poll = ?", true).Error
	if err != nil {
		utils.LogFatal(err.Error())
	}

	// start scheduler
	scheduler.InitScheduler()
	scheduler.StartIntradayPolling(&stocks, db)

	http.ListenAndServe(":8080", stockHandler)
}
