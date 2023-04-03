package main

import (
	"net/http"
	"goTwinder/src/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.RootHandler)
	r.HandleFunc("/swipes/{leftorright}", handlers.SwipesHandler)

	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}

