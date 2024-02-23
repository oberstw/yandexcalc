package handlers

import (
	"net/http"
	"io"
	"time"
	"bytes"
	"sync"
	"strings"
	"encoding/json"
	"orch/math"
)

type JobInfo struct {
	Lock      sync.Mutex     `json:"-"`
	Running   map[string]Job `json:"running"`
	Failed    map[string]Job `json:"failed"`
	Completed map[string]Job `json:"completed"`
}

type Job struct{
	Expr string `json:"expr"`
	Start string `json:"start"`
	End string `json:"end"`
}

type Expres struct{
	Expr string `json:"expr"`
	Id string `json:"id"`
}

type TimeStruct struct{
	Oper string `json:"oper"`
	Timeout int `json:"timeout"`
}

func spaces(line string) string{
	return strings.ReplaceAll(line, " ", "")
}

var JobsTotal JobInfo
var Time = make(map[string]int)
var mu   sync.Mutex 

func SetTimeouts() {
	Time["+"] = time.Duration(3000)
	Time["-"] = time.Duration(3000)
	Time["/"] = time.Duration(3000)
	Time["*"] = time.Duration(3000)
}

func Expr(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var e string
	er := decoder.Decode(&e)
	e = spaces(e)
	if er != nil {
		msg, _ := json.Marshal(0)
		w.Write(msg)
		return
	}
	res := math.InfixToPostfix(e)
	ans, err := math.Calculate(res)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		msg, _ := json.Marshal(0)
		w.Write(msg)
		return
	}
	w.WriteHeader(http.StatusOK)
	msg, _ := json.Marshal(ans)
	w.Write(msg)
}

type Status struct {
	http.ResponseWriter
	Stat int
}

func StartJobs() {
	JobsTotal = JobInfo{sync.Mutex{}, make(map[string]Job), make(map[string]Job), make(map[string]Job)}
}

func Jobhandle(w http.ResponseWriter, r *http.Request) {
	msg, _ := json.Marshal(JobsTotal)
	w.Write(msg)
}

func JobMux(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &Status{ResponseWriter: w}
		body, err := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		data := Expres{}
		err = json.Unmarshal(body, &data)
		if err == nil {
			JobsTotal.Lock.Lock()
			defer JobsTotal.Lock.Unlock()
			JobsTotal.Running[data.Id] = Job{data.Expr, time.Now().Format("2006-01-02 15:04:05"), ""}
		}
		next.ServeHTTP(rec, r)
		if rec.Stat == http.StatusOK {
			JobsTotal.Lock.Lock()
			defer JobsTotal.Lock.Unlock()
			job := JobsTotal.Running[data.Id]
			delete(JobsTotal.Running, data.Id)
			JobsTotal.Completed[data.Id] = Job{job.Expr, job.Start, time.Now().Format("2006-01-02 15:04:05")}
		} else {
			JobsTotal.Lock.Lock()
			defer JobsTotal.Lock.Unlock()
			job := JobsTotal.Running[data.Id]
			delete(JobsTotal.Running, data.Id)
			JobsTotal.Failed[data.Id] = Job{job.Expr, job.Start, time.Now().Format("2006-01-02 15:04:05")}
		}
	})
}

func ChangeExprTimeout(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var e TimeStruct
	decoder.Decode(&e)
	mu.Lock()
	switch e.Oper {
    case "+":
        Time[e.Oper] = time.Duration(e.Timeout)
    case "-":
        Time[e.Oper] = time.Duration(e.Timeout)
    case "/":
        Time[e.Oper] = time.Duration(e.Timeout)
    case "*":
        Time[e.Oper] = time.Duration(e.Timeout)
    default:
        http.Error(w, "Unsupported operation", http.StatusBadRequest)
        return
    }
}

func GibTimeouts(w http.ResponseWriter, r *http.Request) {
	msg, _ := json.Marshal(Time)
	w.Write(msg)
}