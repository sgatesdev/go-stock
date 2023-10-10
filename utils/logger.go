package utils

import (
	"fmt"
	"os"
	"time"
)

func LogMsg(m string) {
	zone, err := time.LoadLocation(os.Getenv("TIMEZONE"))
	if err != nil {
		fmt.Println(err.Error())
	}

	adjusted := time.Now().In(zone)

	fmt.Println(time.Now().Format(time.RFC850) + ": adjusted time: " + adjusted.Format(time.RFC850) + " : " + m)
}

func LogFatal(m string) {
	LogMsg(m)
	os.Exit(1)
}

func LogError(m string) {
	LogMsg(m)
	increaseErrorCount()
}
