package handlers

import (
	"net/http"
	"strconv"
	"encoding/json"
	"io"
	"goTwinder/src/schemas"
	"goTwinder/src/middlewares"
	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
)

var (
	maxNumRandomLiveUsers = 3
	retries = 10
)

func RandomLiveUsersHandler(w http.ResponseWriter, r *http.Request, c schemas.ConnectionCollection) {
	if (r.Method == http.MethodGet) {
		GetRandomLiveUsers(w, r, c)
	} else {
		http.Error(w, "Unsupported request method", http.StatusBadRequest)
	}
}

func GetRandomLiveUsers(w http.ResponseWriter, r *http.Request, c schemas.ConnectionCollection) {
	uid, ok, msg := isRandomUsersUrlValid(r)
	if !ok {
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	middlewares.RefreshUserCache(uid, c.RedisClient)
	retrieved := map[string]struct{}{}

	for i := 0; i < retries && len(retrieved) <= maxNumRandomLiveUsers; i++ {
		result, _ := c.RedisClient.Do("RANDOMKEY").Result()
		if result != nil {
			key, _ := result.(string)
			if key != strconv.Itoa(uid) {
				retrieved[key] = struct{}{}
			}
		}
	}
	stringSlice := make([]string, 0, len(retrieved))
	for k := range retrieved {
		stringSlice = append(stringSlice, k)
	}

	retrievedJSON, _ := json.Marshal(stringSlice)
	io.WriteString(w, string(retrievedJSON))
}

func isRandomUsersUrlValid(r *http.Request) (int, bool, string) {
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