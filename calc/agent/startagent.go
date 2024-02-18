package main

import (
	"net/http"
	"fmt"
	"agent/workers"
	"agent/handlers"
)

func StartAgent() {
	mux := http.NewServeMux()
	exphandler := http.HandlerFunc(exp)
	work := http.HandlerFunc(workers)
	mux.Handle("/exp", exphandler)
	mux.Handle("/workers", work)
	http.ListenAndServe(":8080", mux)
}