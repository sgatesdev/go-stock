package finnhub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"samgates.io/go-stock/models"
	"samgates.io/go-stock/stream"
	"samgates.io/go-stock/utils"
)

var token string

const maxThreads = 5

func PollFinnhub(stocks *[]models.Stock, checkDate func() (bool, string), db *gorm.DB) {
	_, p := os.LookupEnv("DATE_OVERRIDE")

	// get the type of price to fetch
	ok, t := checkDate()
	if !p && !ok {
		return
	}

	var present bool
	token, present = os.LookupEnv("FINNHUB_TOKEN")
	if !present {
		utils.LogFatal("No Finnhub token found - cannot start.")
	}

	// send requests as literals vs. pointers, so we can recursively call
	sendRequests(*stocks, db, t)
}

func sendRequests(stocks []models.Stock, db *gorm.DB, t string) {
	var wg sync.WaitGroup
	var i int
	for n, s := range stocks {
		if n >= maxThreads {
			// ignore
			continue
		} else {
			// count requests
			i++
		}

		wg.Add(1)
		go func(innerS models.Stock, innerT string) {
			defer wg.Done()
			fetchStockQuote(innerS, db, innerT)
		}(s, t)
	}

	wg.Wait()

	if len(stocks) == i {
		// done
		return
	} else {
		// recursively call until we have processed all requests
		time.Sleep(500 * time.Millisecond)
		sendRequests(stocks[i:], db, t)
		return
	}
}

func fetchStockQuote(s models.Stock, db *gorm.DB, t string) {
	c := http.Client{
		Timeout: 5 * time.Second,
	}

	res, err := c.Get("https://finnhub.io/api/v1/quote?symbol=" + s.Symbol + "&token=" + token)
	if err != nil {
		utils.LogError("Error constructing GET request")
		return
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		m := fmt.Sprintf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		utils.LogError(m)
		return
	}
	if err != nil {
		utils.LogError(err.Error())
		return
	}

	p := models.Price{}
	err = json.Unmarshal(body, &p)
	if err != nil {
		utils.LogError(err.Error())
		return
	}

	utils.LogMsg(fmt.Sprintf("Fetched price for %s", s.Symbol))

	stock := &models.Stock{}
	err = db.Model(models.Stock{}).Where("symbol = ?", s.Symbol).First(stock).Error
	if err != nil {
		utils.LogError(err.Error())
		return
	}

	// update model
	p.ID = uuid.NewString()
	p.StockID = stock.ID
	p.Type = t

	// make sure we don't have a duplicate
	dbResponse := db.Where("stock_id = ? AND received = ? AND type = ?", p.StockID, p.Received, t).Find(&models.Price{})
	if dbResponse.Error != nil && dbResponse.Error != gorm.ErrRecordNotFound {
		utils.LogError(dbResponse.Error.Error())
		return
	} else if dbResponse.RowsAffected > 0 {
		// duplicate
		msg := fmt.Sprintf("Duplicate %s price for %s at %d", t, p.StockID, p.Received)
		utils.LogMsg(msg)
		return
	}

	// save
	db.Save(&p)

	// broadcast update to WS connections
	stream.SendPriceUpdate(p)
}
