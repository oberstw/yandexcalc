package handlers

import (
	"net/http"
	"fmt"
	"agent/workers"
	"time"
	"context"
	"encoding/json"
)

type Agentreq struct {
	A     float64 `json:"a"`
	B     float64 `json:"b"`
	Sign    string  `json:"sign"`
}

func RelWorker(wdata string) {
	workers.Limit.Release(1)
	delete(workers.Information, wdata)
}

func Exp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var e Agentreq
	err := decoder.Decode(&e)
	if err != nil {
		panic(err)
	}
	workerdata := fmt.Sprintf("%f %s %f", e.A, e.Sign, e.B)
	workers.Limit.Acquire(context.Background(), 1)
	start := time.Now().Format("2006-01-02 15:04:05")
	workers.Information[workerdata] = start
	defer RelWorker(workerdata)
	var res float64
	var er error
	if e.Sign == "+" {
		res = e.A + e.B
	} else if e.Sign == "-" {
		res = e.A - e.B
	} else if e.Sign == "*" {
		res = e.A * e.B
	} else if e.Sign == "/" {
		if e.B == 0 {
			er = fmt.Errorf("div by zero")
		} else {
			res = e.A /e.B
		}
	} else {
		er = fmt.Errorf("wrong sign")
	}
	time.Sleep(time.Duration(1000) * time.Millisecond)
	if er != nil {
		w.WriteHeader(http.StatusForbidden)
		panic(er)
	}
	k, _ := json.Marshal(res)
	w.Write(k)
}

func Workersh(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(workers.Information)
	w.Write(data)
}