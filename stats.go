package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/manucorporat/stats"
)

var (
	ips        = stats.New()
	mutexStats sync.RWMutex
	savedStats map[string]uint64
)

func statsWorker() {
	c := time.Tick(1 * time.Second)
	var lastMallocs uint64
	var lastFrees uint64

	var x = 0
	for range c {
		var stats runtime.MemStats
		runtime.ReadMemStats(&stats)

		mutexStats.Lock()
		savedStats = map[string]uint64{
			"timestamp":  uint64(time.Now().Unix()),
			"HeapInuse":  stats.HeapInuse,
			"StackInuse": stats.StackInuse,
			"Mallocs":    stats.Mallocs - lastMallocs,
			"Frees":      stats.Frees - lastFrees,
		}
		lastMallocs = stats.Mallocs
		lastFrees = stats.Frees
		mutexStats.Unlock()

		// reset ip counter every 60 seconds if any data is stored
		if len(ips.Data()) > 0 {
			if x == 60 {
				x = 0
				fmt.Println("rate limit: reset ip data")
				ips.Reset()
			} else {
				x++
			}
		}
	}
}

// Stats returns savedStats data.
func Stats() map[string]uint64 {
	mutexStats.RLock()
	defer mutexStats.RUnlock()

	return savedStats
}
