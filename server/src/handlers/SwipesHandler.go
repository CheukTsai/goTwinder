package handlers

import (
	"fmt"
	"io"
	"net/http"
	"encoding/json"
	"goTwinder/src/schemas"
	"github.com/gorilla/mux"
	"strconv"
)

func SwipesHandler(w http.ResponseWriter, r *http.Request) {
	if (r.Method == http.MethodPost) {
		PostSwipes(w, r)
	} else {
		http.Error(w, "Unsupported request method", http.StatusBadRequest)
	}
}

func PostSwipes(w http.ResponseWriter, r *http.Request) {
	isvalid, msg := isUrlValid(r)
	if !isvalid {
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	var swipe schemas.Swipe

	e := json.NewDecoder(r.Body).Decode(&swipe)

	if (e != nil) {
		http.Error(w, "Error decoding JSON request body", http.StatusBadRequest)
        return
	}
	defer r.Body.Close()

	if (!isValidSwipe(&swipe)) {
		http.Error(w, "Missing body attribute", http.StatusBadRequest)
		return
	}

	fmt.Printf("got / POST swipes request\n")
	io.WriteString(w, "Swiper: " + strconv.Itoa(swipe.Swiper))
}

func isValidSwipe(swipe *schemas.Swipe) bool {
	return swipe.Swiper != 0 && swipe.Swipee != 0 && swipe.Comment != ""
}

func isUrlValid(r *http.Request) (bool, string) {
	vars := mux.Vars(r)

	leftOrRight, ok := vars["leftorright"]

	if !ok {
		return false, "Missing parameter leftOrRight"
	}

	if leftOrRight != "left" && leftOrRight != "right" {
		return false, "Wrong parameter, shall be chosen from left and right"
	}

	return true, ""
}