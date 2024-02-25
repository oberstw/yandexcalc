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

	mux.Handle("/expr",  corsMiddleware(handlers.JobMux(expr)))
	mux.Handle("/jobs", corsMiddleware(jobs))
	mux.Handle("/time", corsMiddleware(time))
	mux.Handle("/timeouts", corsMiddleware(givtime))

	http.ListenAndServe(":8000", mux)
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
