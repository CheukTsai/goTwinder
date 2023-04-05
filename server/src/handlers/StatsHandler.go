package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"goTwinder/src/schemas"
	"goTwinder/src/middlewares"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func StatsHandler(w http.ResponseWriter, r *http.Request, c schemas.ConnectionCollection) {
	if (r.Method == http.MethodGet) {
		GetStats(w, r, c)
	} else {
		http.Error(w, "Unsupported request method", http.StatusBadRequest)
	}
}

func GetStats(w http.ResponseWriter, r *http.Request, c schemas.ConnectionCollection) {
	log.Printf("got / GET matches request\n")
	uid, ok, msg := isStatsUrlValid(r)
	if !ok {
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	middlewares.RefreshUserCache(uid, c.RedisClient)
	stmt, err := c.MySqlDatabase.Prepare("SELECT\n" +
    "(SELECT COUNT(*) FROM likes WHERE userid = ?) as num_likes,\n" +
    "(SELECT COUNT(*) FROM dislikes WHERE userid = ?) as num_dislikes")
	
	if err != nil {
        panic(err.Error())
    }
    defer stmt.Close()

	var stats schemas.Stats
	err = stmt.QueryRow(uid, uid).Scan(&stats.NumLikes, &stats.NumDislikes)
	if err != nil {
        panic(err.Error())
    }

	statsJSON, err := json.Marshal(stats)

	if err != nil {
		log.Printf(err.Error())
		http.Error(w, "error parsing", http.StatusBadRequest)
		return
	}

	io.WriteString(w, string(statsJSON))
}

func isStatsUrlValid(r *http.Request) (int, bool, string) {
	vars := mux.Vars(r)

	swiperid, ok := vars["swiperid"]

	if !ok {
		return 0, false, "Missing parameter swiperid"
	}

	id, err := strconv.Atoi(swiperid)

	if err != nil {
		return 0, false, "SwiperId is not a valid number"
	}

	return id, true, ""
}

