package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
	"gorm.io/gorm"
	finnhub "samgates.io/go-stock/finnhub"
	models "samgates.io/go-stock/models"
	utils "samgates.io/go-stock/utils"
)

var Scheduler *gocron.Scheduler

func InitScheduler() {
	Scheduler = gocron.NewScheduler(time.UTC)
}

// starts
func StartIntradayPolling(stocks *[]models.Stock, db *gorm.DB) {
	// _, err := Scheduler.Every(30).Seconds().Tag("intraday").Do(finnhub.PollFinnhub, stocks, CheckDate, db)
	go func() {
		for {
			finnhub.PollFinnhub(stocks, CheckDate, db)
			time.Sleep(30 * time.Second)
		}
	}()
}

func CheckDate() bool {
	t := time.Now()

	y := t.Year()
	m := t.Month()
	weekday := t.Weekday()
	minutes := t.Minute()
	day := t.Day()
	h := t.Hour()

	if weekday == time.Saturday || weekday == time.Sunday {
		// StopIntradayPolling()
		utils.LogMsg("Market closed (weekend)")
		return false
	} else if h < 9 || (h == 9 && minutes < 30) || h >= 16 {
		utils.LogMsg("Market closed (hours)")
		return false
	} else {
		for _, d := range *closeDates() {
			if y == d.year && m == d.month && day == d.day && h >= d.hour {
				utils.LogMsg("Market closed (holiday)")
				return false
			}
		}
	}
	return true
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
