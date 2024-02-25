package main

import (
    "net/http"
)

func main(){
	mux := http.NewServeMux()

	expr := http.HandlerFunc(Expr)
	time := http.HandlerFunc(Time)
	jobs := http.HandlerFunc(Jobs)
	
	mux.Handle("/expr", corsMiddleware(expr))
	mux.Handle("/time", corsMiddleware(time))
	mux.Handle("/jobs", corsMiddleware(jobs))

	http.ListenAndServe(":8040", mux)
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

func Expr(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/main_expr.html")
}

func Jobs(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/jobs.html")
}

func Time(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/time.html")
}