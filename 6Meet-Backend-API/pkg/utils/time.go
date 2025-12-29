package utils

import "time"

// ToDuration converts various number types (seconds) to time.Duration
func ToDuration[T int | int64 | uint64](seconds T) time.Duration {
	return time.Duration(seconds) * time.Second
}

// ToDurationMs converts various number types (milliseconds) to time.Duration
func ToDurationMs[T int | int64 | uint64](milliseconds T) time.Duration {
	return time.Duration(milliseconds) * time.Millisecond
}
