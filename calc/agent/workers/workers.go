package workers

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"
	"golang.org/x/sync/semaphore"
)


var Information map[string]string
var Workers = 10
var Limit semaphore.Weighted


func Set() {
	Limit = *semaphore.NewWeighted(Workers)
	Information = make(map[string]string)
}
