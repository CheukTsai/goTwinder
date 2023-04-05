package main

import (
	"net/http"
	"goTwinder/src/schemas"
	"goTwinder/src/handlers"
	"goTwinder/src/tools"
	"github.com/gorilla/mux"
	"goTwinder/src/managers"
)

var (
	numChannel = 10
)

func initConnections() schemas.ConnectionCollection {
	return schemas.ConnectionCollection{
		RMQChannelPool: tools.NewChannelPool(numChannel),
		MySqlDatabase:  managers.MySqlConnectDatabase(),
		RedisClient:  managers.NewRedisClientConnection(numChannel, numChannel),
	}
}

func main() {
	connections := initConnections()
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.RootHandler)
	r.HandleFunc("/swipes/{leftorright}", func(w http.ResponseWriter, r *http.Request) {
		handlers.SwipesHandler(w, r, connections)
	})
	r.HandleFunc("/matches/{swiperid}", func(w http.ResponseWriter, r *http.Request) {
		handlers.MatchesHandler(w, r, connections)
	})
	r.HandleFunc("/stats/{swiperid}", func(w http.ResponseWriter, r *http.Request) {
		handlers.StatsHandler(w, r, connections)
	})
	r.HandleFunc("/random-live-users/{swiperid}", func(w http.ResponseWriter, r *http.Request) {
		handlers.RandomLiveUsersHandler(w, r, connections)
	})
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}
