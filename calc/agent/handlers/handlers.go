package handlers

import (
	"net/http"
	"fmt"
	"orch/math"
	"agent/workers"
	"time"
	"sync"
)

func Exp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var e math.Agentreq
	err := decoder.Decode(&e)
	if err != nil {
		panic(err)
	}
	workerdata := fmt.Sprintf("%f %f %s", data.A, data.B, data.Sign)
	workers.Limit.Acquire(context.Background(), 1)
	start := time.Now().Format("2006-01-02 15:04:05")
	workers.Information[workerdata] = start
	defer func {
		workers.Limit.Release(1)
		delete(workers.Information, workerdata)
	}
	if e.Sign == "+" {
		res := e.A + e.B
	} else if e.Sign == "-" {
		res := e.A - e.B
	} else if e.Sign == "*" {
		res := e.A * e.B
	} else if e.Sign == "/" {
		if e.B == 0 {
			err := fmt.Errorf("div by zero")
		} else {
			res := e.A /e.B
		}
	} else {
		err := fmt.Errorf("wrong sign")
	}
	time.Sleep(time.Duration(1000) * time.Millisecond)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		panic(err)
	}
	k, _ := json.Marshal(res)
	w.Write(res)
}
