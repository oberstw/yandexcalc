package agent

import (
	"net/http"
	"fmt"
	"agent/workers"
	"time"
)

type req struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
	Sign string `json:"sign"`
	Timeout time.Time() `json:"timeout"`
}

func (r *req) evex() float64, error {
	if r.Sign == "+" {
		res := r.A + r.B
	} else if r.Sign == "-" {
		res := r.A - r.B
	} else if r.Sign == "*" {
		res := r.A + r.B
	} else if r.Sign == "/" {
		if r.B == 0 {
			return 0, fmt.Errorf("Division by 0")
		}
		res := r.A / r.B
	} else {
		return 0, fmt.Errorf("Wrong sign")
	}
	time.Sleep(time.Duration(r.Timeout) * time.Millisecond)
	return res, nil
}

func Exp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var e req
	err := decoder.Decode(&e)
	if err != nil {
		panic(err)
	}
	res, err := e.evex()
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		panic(err)
	}
	k, _ := json.Marshal(res)
	w.Write(res)
}

func Workers(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(&workers.Information)
	w.Write(data)
}
