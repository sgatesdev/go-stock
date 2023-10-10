package handlers

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	models "samgates.io/go-stock/models"
)

// get prices by stock id
type PriceHandler struct {
	Router *mux.Router
	db     *gorm.DB
}

// NewStockHandler creates a new stock handler
func NewPriceHandler(router *mux.Router, db *gorm.DB) *PriceHandler {
	return &PriceHandler{
		Router: router,
		db:     db,
	}
}

// ServeHTTP implements the http.Handler interface
func (h *PriceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Router.ServeHTTP(w, r)
}

// RegisterRoutes registers routes for the stock handler
func (h *PriceHandler) RegisterRoutes() {
	h.Router.HandleFunc("/prices/{stockId}", h.handleGetPrice).Methods("GET")
}

// handleGetStock handles getting a stock
func (h *PriceHandler) handleGetPrice(w http.ResponseWriter, r *http.Request) {
	stockId := mux.Vars(r)["stockId"]

	prices := []models.Price{}
	res := h.db.Where("stock_id = ?", stockId).Find(&prices)

	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	pLite := transformPrices(prices)
	json.NewEncoder(w).Encode(&pLite)
}

type StockPrice struct {
	Received int     `json:"t"`
	Type     string  `json:"type"`
	Price    float32 `json:"c"`
}

func transformPrices(prices []models.Price) []StockPrice {
	sort.SliceStable(prices, func(a, b int) bool {
		return prices[a].CreatedAt.UnixMilli() < prices[b].CreatedAt.UnixMilli()
	})

	res := []StockPrice{}
	for _, p := range prices {
		res = append(res, StockPrice{
			Price:    p.Price,
			Received: p.Received,
			Type:     p.Type,
		})
	}

	return res
}
