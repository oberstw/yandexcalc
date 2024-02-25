package workers

import (
	"golang.org/x/sync/semaphore"
)


var Information map[string]string
var Workers = 10
var Limit semaphore.Weighted


func Set() {
	Limit = *semaphore.NewWeighted(int64(Workers))
	Information = make(map[string]string)
}
