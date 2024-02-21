package starter

import (
	"net/http"
	"agent/handlers"
)

func StartAgent() {
	mux := http.NewServeMux()
	exphandler := http.HandlerFunc(handlers.Exp)
	mux.Handle("/exp", exphandler)
	http.ListenAndServe(":8080", mux)
}