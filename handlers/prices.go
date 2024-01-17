package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	models "samgates.io/go-stock/models"
	"samgates.io/go-stock/stream"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

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
	h.Router.HandleFunc("/ws/prices", h.handleStreamPrices).Methods("GET")
}

// handleGetStock handles getting a stock
func (h *PriceHandler) handleGetPrice(w http.ResponseWriter, r *http.Request) {
	stockId := mux.Vars(r)["stockId"]

	qp := r.URL.Query()
	today := qp.Get("today")

	prices := []models.Price{}

	where := "stock_id = ?"
	whereArgs := []interface{}{stockId}

	var res *gorm.DB
	if today == "true" {
		where += " AND DATE(created_at) = ?"
		t := time.Now().Format("2006-01-02")
		whereArgs = append(whereArgs, t)
	}

	res = h.db.Where(where, whereArgs...).Find(&prices)

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

// router.HandleFunc("/echo", echo)
func (h *PriceHandler) handleStreamPrices(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	// register session so polling will broadcast updates
	sessionId := stream.AddConnection(c)

	defer func() {
		c.Close()
		stream.RemoveConnection(sessionId)
	}()

	// testing
	websocketTest, ok := os.LookupEnv("WEBSOCKET_TEST")
	if ok && websocketTest == "true" {
		fmt.Println("WEBSOCKET_TEST")
		go func() {
			for {
				rand.Seed(time.Now().UnixNano())

				// test to see if web sockets will broadcast updates
				time.Sleep(3 * time.Second)
				stocks := []models.Stock{}
				h.db.Find(&stocks).Where("poll = ?", true)
				for _, s := range stocks {
					lastPrice := models.Price{}
					h.db.Where("stock_id = ? AND type = ?", s.ID, "intraday").Order("created_at desc").First(&lastPrice)
					r := float32(math.Round(float64(lastPrice.Price+(rand.Float32()-0.5)*2)*100) / 100)
					fmt.Println(lastPrice.Price, r)
					p := models.Price{
						ID:      uuid.NewString(),
						StockID: s.ID,
						Price:   r,
					}
					stream.SendPriceUpdate(p)
				}
			}
		}()
	}

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
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
