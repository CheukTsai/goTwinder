package handlers

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	fmt.Printf("got / request\n")
	fmt.Printf("%s\n", string(body))
	io.WriteString(w, string(body))
}