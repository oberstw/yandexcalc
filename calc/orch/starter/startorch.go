package starter

import (
	"net/http"
	"orch/handlers"
)

func StartOrchestrator(){
	handlers.StartJobs()
	mux := http.NewServeMux()
	expr := http.HandlerFunc(handlers.Expr)
	jobs := http.HandlerFunc(handlers.Jobhandle)

	mux.Handle("/expr",  handlers.JobMux(expr))
	mux.HandleFunc("/jobs", jobs)

	http.ListenAndServe("8000", mux)
}