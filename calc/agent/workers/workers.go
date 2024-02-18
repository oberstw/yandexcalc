package workers

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"
	"golang.org/x/sync/semaphore"
)

type WorkersInfo struct {
	workers int64              `json:"-"`
	limit   semaphore.Weighted `json:"-"`
	Current map[string]string  `json:"current"`
	lock    sync.Mutex         `json:"-"`
}

var Information WorkersInfo

func Set() {
	env := os.Getenv("MAX_WORKERS")
	workers, err := strconv.ParseInt(env, 10, 64)
	if err != nil {
		logging.ReportAction("did not find env MAX_WORKERS, setting default 10")
		workers = 10
	}
	Information = WorkersInfo{workers, *semaphore.NewWeighted(workers), make(map[string]string), sync.Mutex{}}
}