package starter

import (
	"net/http"
	"agent/handlers"
)

func StartAgent() {
	mux := http.NewServeMux()
	exphandler := http.HandlerFunc(handlers.Exp)
	workhandler := http.HandlerFunc(handlers.Workersh)
	mux.Handle("/exp", corsMiddleware(exphandler))
	mux.Handle("/workers", corsMiddleware(workhandler))
	http.ListenAndServe(":8080", mux)
}

func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}
