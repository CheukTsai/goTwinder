package handlers

import (
	"log"
	"net/http"
	"strconv"
	"fmt"
	"encoding/json"
	"io"
	"goTwinder/src/schemas"
	"goTwinder/src/middlewares"
	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
)

func MatchesHandler(w http.ResponseWriter, r *http.Request, c schemas.ConnectionCollection) {
	if (r.Method == http.MethodGet) {
		GetMatches(w, r, c)
	} else {
		http.Error(w, "Unsupported request method", http.StatusBadRequest)
	}
}

func GetMatches(w http.ResponseWriter, r *http.Request, c schemas.ConnectionCollection) {
	log.Printf("got / GET matches request\n")
	uid, ok, msg := isMatchesUrlValid(r)
	if !ok {
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	middlewares.RefreshUserCache(uid, c.RedisClient)
	stmt, err := c.MySqlDatabase.Prepare("SELECT DISTINCT l1.swipeeid FROM likes l1\n" + 
	"JOIN likes l2 ON l1.swipeeid = l2.userid AND l1.userid = l2.swipeeid\n" +
	"WHERE l1.userid = ?")
	
	if err != nil {
        panic(err.Error())
    }
    defer stmt.Close()
	rows, err := stmt.Query(uid)
	if err != nil {
        panic(err.Error())
    }
    defer rows.Close()
	matches := []string{}
	for rows.Next() {
        var userid string
        err = rows.Scan(&userid)
        if err != nil {
            panic(err.Error())
        }
		matches = append(matches, userid)
        fmt.Printf("UserID: %s\n", userid)
    }
	matchesJSON, _ := json.Marshal(matches)
	io.WriteString(w, string(matchesJSON))
}

func isMatchesUrlValid(r *http.Request) (int, bool, string) {
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

