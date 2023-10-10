package handlers

// create new stock
// update stock
// delete stock
// get stock
// get all stocks
// get all stocks for a user

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"

	models "samgates.io/go-stock/models"
)

// StockHandler handles stocks
type StockHandler struct {
	Router *mux.Router
	db     *gorm.DB
}

// NewStockHandler creates a new stock handler
func NewStockHandler(router *mux.Router, db *gorm.DB) *StockHandler {
	return &StockHandler{
		Router: router,
		db:     db,
	}
}

// ServeHTTP implements the http.Handler interface
func (h *StockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Router.ServeHTTP(w, r)
}

// RegisterRoutes registers routes for the stock handler
func (h *StockHandler) RegisterRoutes() {
	h.Router.HandleFunc("/stocks/", h.handleGetAllStocks).Methods("GET")
	h.Router.HandleFunc("/stocks/{id}", h.handleGetStock).Methods("GET")
	h.Router.HandleFunc("/stocks/one", h.handleCreateStock).Methods("POST")
	h.Router.HandleFunc("/stocks/many", h.handleCreateStocks).Methods("POST")
	h.Router.HandleFunc("/stocks/{id}", h.handleUpdateStock).Methods("PUT")
	h.Router.HandleFunc("/stocks/{id}", h.handleDeleteStock).Methods("DELETE")
}

// handleGetAllStocks handles getting all stocks
func (h *StockHandler) handleGetAllStocks(w http.ResponseWriter, r *http.Request) {
	stocks := []models.Stock{}
	err := h.db.Find(&stocks).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	collection := models.Stocks{Stocks: stocks}
	json.NewEncoder(w).Encode(collection)
}

// handleGetStock handles getting a stock
func (h *StockHandler) handleGetStock(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	stock := models.Stock{}
	res := h.db.Where("ID = ?", id).First(&stock)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stock)
}

// handleCreateStock handles creating a stock
func (h *StockHandler) handleCreateStock(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	stock := models.Stock{}
	err = json.Unmarshal(body, &stock)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// create stock
	// save stock
	// return stock
	stock.ID = uuid.NewString()
	h.db.Save(&stock)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stock)
}

// handleCreateStock handles creating several stocks
func (h *StockHandler) handleCreateStocks(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	stocks := models.Stocks{}
	err = json.Unmarshal(body, &stocks)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	collection := []models.Stock{}
	for _, s := range stocks.Stocks {
		s.ID = uuid.NewString()
		h.db.Save(&s)
		collection = append(collection, s)
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(collection)
}

// handleUpdateStock handles updating a stock
func (h *StockHandler) handleUpdateStock(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	stock := models.Stock{}
	err = json.Unmarshal(body, &stock)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := h.db.Where("ID = ?", id).First(&models.Stock{})
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	stock.ID = id
	err = h.db.Save(&stock).Error
	fmt.Println(stock)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stock)
}

// handleDeleteStock handles deleting a stock
func (h *StockHandler) handleDeleteStock(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := h.db.Where("ID = ?", id).First(&models.Stock{})
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := h.db.Where("stock_id = ?", id).Delete(&models.Price{}).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s := &models.Stock{ID: id}
	err = h.db.Delete(&s).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// getStocks gets all stocks
// func getStocks() ([]Stock, error) {
// }

// // getStock gets a stock
// func getStock(id int) (Stock, error) {
// }

// // createStock creates a stock
// func createStock(s Stock) (Stock, error) {
// }

// // updateStock updates a stock
// func updateStock(s Stock) (Stock, error) {
// }

// // deleteStock deletes a stock
// func deleteStock(id int) error {
// }
