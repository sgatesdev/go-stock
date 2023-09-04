package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

func LogMsg(m string) {
	var (
		buf    bytes.Buffer
		logger = log.New(&buf, "go-stock: ", log.Ldate|log.Ltime)
	)

	logger.Print(m)
	fmt.Print(&buf)
}

func LogFatal(m string) {
	var (
		buf    bytes.Buffer
		logger = log.New(&buf, "go-stock: ", log.Ldate|log.Ltime)
	)

	logger.Print(m)
	fmt.Print(&buf)
	os.Exit(1)
}

func LogError(m string) {
	var (
		buf    bytes.Buffer
		logger = log.New(&buf, "go-stock: ", log.Ldate|log.Ltime)
	)

	logger.Print(m)
	fmt.Print(&buf)

	increaseErrorCount()
}
