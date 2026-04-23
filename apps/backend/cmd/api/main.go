package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	println("Server starting on 8080")
	http.ListenAndServe(":8080", r)
}
