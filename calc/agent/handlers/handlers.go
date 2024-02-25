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
	Timeout int `json:"timeout"`
}

type Agentout struct {
	Result float64 `json:"result"`
	Err string  `json:"err"`
}

type WorkData struct {
	Data map[string]string `json:"workdata"`
}

func RelWorker(wdata string) {
	workers.Limit.Release(1)
	delete(workers.Information, wdata)
}

func Exp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Does this work?")
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
	fmt.Println("Worker opened", workers.Information)
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
	time.Sleep(time.Duration(e.Timeout) * time.Millisecond)
	if er != nil {
		w.WriteHeader(http.StatusForbidden)
		panic(er)
	}
	fmt.Println("Agent thing done successfully", res)
	var g Agentout
	g.Result = res
	k, _ := json.Marshal(g)
	w.Write(k)
}

func Workersh(w http.ResponseWriter, r *http.Request) {
	wd, _ := json.Marshal(workers.Information)
	fmt.Println("Workers returned ", wd)
	w.Write(wd)
}