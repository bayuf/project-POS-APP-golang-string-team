package utils

import (
	"fmt"
	"sync/atomic"
	"time"
)

var (
	lastDate int64
	counter  uint64
)

func GenerateOrderNum() string {
	now := time.Now()
	today := now.Year()*10000 + int(now.Month())*100 + now.Day()

	// reset counter jika tanggal berubah
	if atomic.LoadInt64(&lastDate) != int64(today) {
		atomic.StoreInt64(&lastDate, int64(today))
		atomic.StoreUint64(&counter, 0)
	}

	n := atomic.AddUint64(&counter, 1)

	return fmt.Sprintf(
		"ORD-%04d%02d%02d-%06d",
		now.Year(),
		now.Month(),
		now.Day(),
		n,
	)
}
