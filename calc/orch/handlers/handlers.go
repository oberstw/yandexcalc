package orch

import (
	"net/http"
	"fmt"
	"orch/math"
)

type JobInfo struct {
	Lock      sync.Mutex     `json:"-"`
	Running   map[string]Job `json:"running"`
	Failed    map[string]Job `json:"failed"`
	Completed map[string]Job `json:"completed"`
}

type Expression struct {
	Id   string `json:"id"`
	Expr string `json:"expr"`
}

type Result struct {
	Res float64 `json:"res"`
	Err string  `json:"err"`
}


type Job struct{
	Expr string `json:"expr"`
	Start string `json:"start"`
	End string `json:"end"`
}

var JobsTotal *JobInfo

func Expr(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var e string
	er := decoder.Decode(&e)
	e = math.sanitize(e)
	if er != nil {
		ermsg := Result{0, er}
		msg, _ := json.Marshal(ermsg)
		w.Write(msg)
		return
	}
	res := math.InfixToPostfix(e)
	ans, err := math.Calculate(res)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		ermsg := Result{0, err}
		msg, _ := json.Marshal(ermsg)
		w.Write(msg)
		return
	}
	w.WriteHeader(http.SatusOK)
	res := Result{ans, ""}
	msg, _ := json.Marshal(res)
	w.Write(msg)
}

type Status struct {
	http.ResponseWriter
	Stat int
}

func Start() {
	Jobs = &JobsInfo{sync.Mutex{}, make(map[string]Job), make(map[string]Job), make(map[string]Job)}
}

func Jobhandle(w http.ResponseWriter, r *http.Request) {
	msg, _ := json.Marshal(JobsTotal)
	w.Write(msg)
}

func JobMux(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &StatusRecorder{ResponseWriter: w}
		rec.Stat = 200
		body, err := io.ReadAll(r.Body)

		r.Body = io.NopCloser(bytes.NewBuffer(body))
		data := AddExprReqIn{}
		err = json.Unmarshal(body, &data)
		if err == nil {
			addJob(data)
		}

		next.ServeHTTP(rec, r)
		if rec.Stat == http.StatusOK {
			toComp(data)
		} else {
			toFail(data)
		}
	})
}

func addJob(data Expression) {
	JobsTotal.Lock.Lock()
	defer JobsTotal.Lock.Unlock()
	JobsTotal.Running[data.Id] = Job{data.Expr, time.Now().Format("01/02 - 03:04:05"), ""}
}

func toComp(data Expression) {
	JobsTotal.Lock.Lock()
	defer JobsTotal.Lock.Unlock()
	job := JobsTotal.Running[data.Id]
	delete(JobsTotal.Running, data.Id)
	JobsTotal.Completed[data.Id] = Job{job.Expr, job.Start, time.Now().Format("01/02 - 03:04:05")}
}

func toFail(data Expression) {
	JobsTotal.Lock.Lock()
	defer JobsTotal.Lock.Unlock()
	job := JobsTotal.Running[data.Id]
	delete(JobsTotal.Running, data.Id)
	JobsTotal.Failed[data.Id] = Job{job.Expr, job.Start, time.Now().Format("01/02 - 03:04:05")}
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}