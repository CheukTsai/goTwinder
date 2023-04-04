package handlers

import (
	"context"
	"log"
	"time"
	"io"
	"net/http"
	"encoding/json"
	"goTwinder/src/schemas"
	"goTwinder/src/tools"
	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
	"strconv"
)

func SwipesHandler(w http.ResponseWriter, r *http.Request, cp *tools.ChannelPool) {
	if (r.Method == http.MethodPost) {
		PostSwipes(w, r, cp)
	} else {
		http.Error(w, "Unsupported request method", http.StatusBadRequest)
	}
}

func PostSwipes(w http.ResponseWriter, r *http.Request, cp *tools.ChannelPool) {
	log.Printf("got / POST swipes request\n")
	isvalid, msg := isUrlValid(r)
	if !isvalid {
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var swipe schemas.Swipe

	e := json.NewDecoder(r.Body).Decode(&swipe)

	if (e != nil) {
		http.Error(w, "Error decoding JSON request body", http.StatusBadRequest)
        return
	}
	

	if (!isValidSwipe(&swipe)) {
		http.Error(w, "Missing body attribute", http.StatusBadRequest)
		return
	}
	swipeJSON, err := json.Marshal(swipe)

	if err != nil {
		log.Printf(err.Error())
		http.Error(w, "error parsing", http.StatusBadRequest)
		return
	}
	ch := cp.Get()
	defer cp.Put(ch)
	sendMsgToRMQ(string(swipeJSON), ch)
	io.WriteString(w, "Swiper: " + strconv.Itoa(swipe.Swiper))
}

func sendMsgToRMQ(msg string, ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		"swipeQueue", // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	  )
	  failOnError(err, "Failed to declare a queue")
	  
	  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	  defer cancel()
	  
	  body := msg
	  err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
		  ContentType: "text/plain",
		  Body:        []byte(body),
		})
	  failOnError(err, "Failed to publish a message")
	  log.Printf(" [x] Sent %s\n", body)
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

func failOnError(err error, msg string) {
	if err != nil {
	  log.Panicf("%s: %s", msg, err)
	}
}