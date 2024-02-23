package starter

import (
	"net/http"
	"orch/handlers"
)

func StartOrchestrator(){
	handlers.StartJobs()
	handlers.SetTimeouts()
	mux := http.NewServeMux()
	expr := http.HandlerFunc(handlers.Expr)
	jobs := http.HandlerFunc(handlers.Jobhandle)
	time := http.HandlerFunc(handlers.ChangeExprTimeout)
	givtime := http.HandlerFunc(handlers.GibTimeouts)

	mux.Handle("/expr",  handlers.JobMux(expr))
	mux.HandleFunc("/jobs", jobs)
	mux.HandleFunc("/time", time)
	mux.HandleFunc("/timeouts", givtime)

	http.ListenAndServe(":8000", mux)
}