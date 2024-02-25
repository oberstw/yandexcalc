package handlers

import (
	"net/http"
	"io"
	"time"
	"bytes"
	"fmt"
	"github.com/urfave/negroni"
	"strings"
	"encoding/json"
	"orch/math"
)

type JobInfo struct {
	Running   map[string]Job `json:"running"`
	Completed map[string]Job `json:"completed"`
	Failed    map[string]Job `json:"failed"`
}

type Job struct{
	Expr string `json:"expr"`
	Ans float64 `json:"ans"`
	Start string `json:"start"`
	End string `json:"end"`
}

type Expres struct{
	Expr string `json:"expr"`
	Id string `json:"id"`
}

type Res struct{
	Rs float64 `json:"result"`
	Err error `json:"err"`
}

type TimeStruct struct{
	Oper string `json:"oper"`
	Timeout int `json:"timeout"`
}

func spaces(line string) string{
	return strings.ReplaceAll(line, " ", "")
}

var JobsTotal JobInfo

func SetTimeouts() {
	math.Time["+"] = 3000
	math.Time["-"] = 3000
	math.Time["/"] = 3000
	math.Time["*"] = 3000
}

func Expr(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var e Expres
	er := decoder.Decode(&e)
	a := spaces(e.Expr)
	fmt.Println(a)
	if er != nil {
		msg, _ := json.Marshal(0)
		w.Write(msg)
		return
	}
	fmt.Println("Before ItoP")
	res := math.InfixToPostfix(a)
	fmt.Println("ItoP successful")
	ans, err := math.Calculate(res)
	fmt.Println("Calculate successful")
	if err != nil {
		fmt.Println("Error found after Calculate")
		w.WriteHeader(http.StatusForbidden)
		msg, _ := json.Marshal(0)
		w.Write(msg)
		return
	}
	w.WriteHeader(http.StatusOK)
	msg, _ := json.Marshal(ans)
	w.Write(msg)
}

func StartJobs() {
	JobsTotal = JobInfo{make(map[string]Job), make(map[string]Job), make(map[string]Job)}
}

func Jobhandle(w http.ResponseWriter, r *http.Request) {
	msg, _ := json.Marshal(JobsTotal)
	fmt.Println("got called by js")
	w.Write(msg)
}

func JobMux(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := negroni.NewResponseWriter(w)
		body, err := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		data := Expres{}
		err = json.Unmarshal(body, &data)
		if err == nil {
			JobsTotal.Running[data.Id] = Job{data.Expr, 0,  time.Now().Format("2006-01-02 15:04:05"), ""}
		}
		next.ServeHTTP(rec, r)
		if rec.Status() == http.StatusOK {
			fmt.Println("Status OK")
			job := JobsTotal.Running[data.Id]
			fmt.Println(job)
			delete(JobsTotal.Running, data.Id)
			JobsTotal.Completed[data.Id] = Job{job.Expr, 0, job.Start, time.Now().Format("2006-01-02 15:04:05")}
		} else {
			fmt.Println("Status not OK")
			job := JobsTotal.Running[data.Id]
			delete(JobsTotal.Running, data.Id)
			JobsTotal.Failed[data.Id] = Job{job.Expr, 0, job.Start, time.Now().Format("2006-01-02 15:04:05")}
			fmt.Println(JobsTotal.Failed[data.Id])
		}
		fmt.Println("Job mux check", JobsTotal)
	})
}

func ChangeExprTimeout(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var e TimeStruct
	decoder.Decode(&e)
	fmt.Println("Does it work?")
	fmt.Println(e)
	if e.Oper == "+"{
		math.Time[e.Oper] = e.Timeout
	} else if e.Oper == "-" {
		math.Time[e.Oper] = e.Timeout
	} else if e.Oper == "/" {
		math.Time[e.Oper] = e.Timeout
	} else if e.Oper == "*" {
		math.Time[e.Oper] = e.Timeout
	} else {
        http.Error(w, "Unsupported operation", http.StatusBadRequest)
        return
	}
	fmt.Println("Timeout set")
}

func GibTimeouts(w http.ResponseWriter, r *http.Request) {
	fmt.Println(math.Time)
	msg, _ := json.Marshal(math.Time)
	w.Write(msg)
}
