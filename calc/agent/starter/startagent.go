package starter

import (
	"net/http"
	"agent/handlers"
)

func StartAgent() {
	mux := http.NewServeMux()
	exphandler := http.HandlerFunc(handlers.Exp)
	workhandler := http.HandlerFunc(handlers.Workersh)
	mux.Handle("/exp", exphandler)
	mux.Handle("/workers", workhandler)
	http.ListenAndServe(":8080", mux)
}