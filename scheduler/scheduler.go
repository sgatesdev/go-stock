package scheduler

import (
	"os"
	"time"

	"gorm.io/gorm"
	finnhub "samgates.io/go-stock/finnhub"
	models "samgates.io/go-stock/models"
	utils "samgates.io/go-stock/utils"
)

// starts
func StartIntradayPolling(stocks *[]models.Stock, db *gorm.DB) {
	go func() {
		for {
			finnhub.PollFinnhub(stocks, CheckDate, db)
			time.Sleep(30 * time.Second)
		}
	}()
}

// returns true and the type of price to fetch
// closing prices are fetched 30 mins before market open
func CheckDate() (bool, string) {
	zone, err := time.LoadLocation(os.Getenv("TIMEZONE"))
	if err != nil {
		utils.LogFatal(err.Error())
	}

	t := time.Now().In(zone)
	y := t.Year()
	m := t.Month()
	weekday := t.Weekday()
	minutes := t.Minute()
	day := t.Day()
	h := t.Hour()

	// DO NOT POLL
	if weekday == time.Saturday || weekday == time.Sunday {
		// StopIntradayPolling()
		utils.LogMsg("Market closed (weekend)")
		return false, ""
	} else if h < 9 || h >= 16 {
		utils.LogMsg("Market closed (hours)")
		return false, ""
	} else {
		for _, d := range *closeDates() {
			if y == d.year && m == d.month && day == d.day && h >= d.hour {
				utils.LogMsg("Market closed (holiday)")
				return false, ""
			}
		}
	}

	// POLL
	if h == 9 && minutes < 30 {
		return true, "last_close"
	}
	return true, "intraday"
}

func closeDates() *[]closeDate {
	// TODO: store in database
	dates := make([]closeDate, 0)
	dates = append(dates, closeDate{year: 2023, month: time.September, day: 4, hour: 9})
	dates = append(dates, closeDate{year: 2023, month: time.November, day: 23, hour: 13})
	dates = append(dates, closeDate{year: 2023, month: time.December, day: 25, hour: 9})

	return &dates
}

type closeDate struct {
	year  int
	month time.Month
	day   int
	hour  int
}
