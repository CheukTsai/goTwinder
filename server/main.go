package main

import (
	"net/http"
	"goTwinder/src/handlers"
	"goTwinder/src/tools"
	"github.com/gorilla/mux"
	"goTwinder/src/managers"
)

var (
	numChannel = 10
)

func main() {
	// db := managers.MySqlConnectDatabase()
	cp := tools.NewChannelPool(20)
	db := managers.MySqlConnectDatabase()
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.RootHandler)
	r.HandleFunc("/swipes/{leftorright}", func(w http.ResponseWriter, r *http.Request) {
		handlers.SwipesHandler(w, r, cp)
	})
	r.HandleFunc("/matches/{swiperid}", func(w http.ResponseWriter, r *http.Request) {
		handlers.MatchesHandler(w, r, db)
	})
	r.HandleFunc("/stats/{swiperid}", func(w http.ResponseWriter, r *http.Request) {
		handlers.StatsHandler(w, r, db)
	})
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}
