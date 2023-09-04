package utils

import "sync"

const maxErrorCount = 5

type ErrorCount struct {
	mu sync.Mutex
	n  int
}

var ECount ErrorCount

func increaseErrorCount() {
	ECount.mu.Lock()
	ECount.n++
	ECount.mu.Unlock()

	if GetErrorCount() > maxErrorCount {
		LogFatal("Error count exceeded")
	}
}

func resetErrorCount() {
	ECount.mu.Lock()
	ECount.n = 0
	ECount.mu.Unlock()
}

func GetErrorCount() int {
	ECount.mu.Lock()
	defer ECount.mu.Unlock()
	return ECount.n
}
