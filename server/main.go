package main

import (
	"net/http"
	"goTwinder/src/handlers"
	"goTwinder/src/tools"
	"github.com/gorilla/mux"
)

func main() {
	cp := tools.NewChannelPool(20)
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.RootHandler)
	r.HandleFunc("/swipes/{leftorright}", func(w http.ResponseWriter, r *http.Request) {
		handlers.SwipesHandler(w, r, cp)
	})

	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}

