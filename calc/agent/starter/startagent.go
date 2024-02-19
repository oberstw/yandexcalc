package main

import (
	"net/http"
	"fmt"
	"agent/handlers"
)

func StartAgent() {
	mux := http.NewServeMux()
	exphandler := http.HandlerFunc(handlers.Exp)
	mux.Handle("/exp", exphandler)
	http.ListenAndServe(":8080", mux)
}