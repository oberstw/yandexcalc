package orch

import (
	"net/http"
	"fmt"
	"orch/handlers"
)

func StartOrchestrator(){
	handlers.Start()
	mux := http.NewServeMux()
	expr := http.HandlerFunc(handlers.Expr)
	jobs := http.HandlerFunc(handlers.Job)

	mux.Handle("/expr", handlers.Middleware(expr))
	mux.HandleFunc("/jobs", handlers.JobMux(jobs))

	http.ListenAndServe("8000", mux)
}