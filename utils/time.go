package utils

import "time"

// GetMillionSecondTimestamp will return current timestamp (ms)
func GetMillionSecondTimestamp() int64 {
	return time.Now().UnixNano() / 1e6
}
