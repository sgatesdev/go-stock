package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	handlers "samgates.io/go-stock/handlers"
	models "samgates.io/go-stock/models"
	scheduler "samgates.io/go-stock/scheduler"
	"samgates.io/go-stock/utils"
)

func main() {
	// connect to db
	db := utils.SetupDB()

	// create router
	router := mux.NewRouter()

	// add stock routes
	stockHandler := handlers.NewStockHandler(router, db)
	stockHandler.RegisterRoutes()

	// add price routes
	priceHandler := handlers.NewPriceHandler(router, db)
	priceHandler.RegisterRoutes()

	// add scheduler routes
	// schedulerHandler := handlers.NewSchedulerHandler(router, db)
	// schedulerHandler.RegisterRoutes()

	// fetch stocks
	stocks := []models.Stock{}
	err := db.Find(&stocks).Where("poll = ?", true).Error
	if err != nil {
		utils.LogFatal(err.Error())
	}

	// start scheduler
	scheduler.StartIntradayPolling(&stocks, db)

	handler := cors.Default().Handler(router)

	http.ListenAndServe(":8080", handler)
}
