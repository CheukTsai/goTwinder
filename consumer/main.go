package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"goTwinderRMQConsumer/src/schemas"
	"goTwinderRMQConsumer/src/tools"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var numChannel = 20
  
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func establishSQLconnection() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	connectionName := fmt.Sprintf("%s:%s@/%s", dbUser, dbPassword, dbName)
	db, err := sql.Open("mysql", connectionName)
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(60)
	db.SetMaxIdleConns(60)
	return db
}

func main() {
	db := establishSQLconnection()
	cp := tools.NewChannelPool(numChannel)
	
	for i := 0; i < numChannel; i++ {
		i := i
		go func() {
			for {
				ch := cp.Get()
				log.Printf(" [*] Waiting for messages. To exit press CTRL+C thread %d", i)
				q, err := ch.QueueDeclare(
					"swipeQueue", // name
					true,   // durable
					false,   // delete when unused
					false,   // exclusive
					false,   // no-wait
					nil,     // arguments
					)
					failOnError(err, "Failed to declare a queue")
				
				msgs, err := ch.Consume(
					q.Name, // queue
					"",     // consumer
					true,   // auto-ack
					false,  // exclusive
					false,  // no-local
					false,  // no-wait
					nil,    // args
					)
				failOnError(err, "Failed to register a consumer")
				for d := range msgs {
					log.Printf("Received a message: %s from thread %d", d.Body, i)
	
					var swipe schemas.Swipe
	
					e := json.NewDecoder(bytes.NewReader(d.Body)).Decode(&swipe)

					if (e != nil) {
						failOnError(e, "Failed to decode json")
					}
					
					if swipe.Like {
						stmt, err := db.Prepare("INSERT INTO likes(userid, swipeeid) VALUES(?, ?)")
						failOnError(err, "Failed on preparing statement")
						defer stmt.Close()
						res, err := stmt.Exec(strconv.Itoa(swipe.Swiper), strconv.Itoa(swipe.Swipee))
						if err != nil {
							log.Fatal(err)
						}
					
						// Print the ID of the newly inserted row
						id, err := res.LastInsertId()
						if err != nil {
							log.Fatal(err)
						}
						log.Printf("Inserted row with ID %d in Likes", id)
					} else {
						stmt, err := db.Prepare("INSERT INTO dislikes(userid, swipeeid) VALUES(?, ?)")
						failOnError(err, "Failed on preparing statement")
						defer stmt.Close()
						res, err := stmt.Exec(strconv.Itoa(swipe.Swiper), strconv.Itoa(swipe.Swipee))
						if err != nil {
							log.Fatal(err)
						}
					
						// Print the ID of the newly inserted row
						id, err := res.LastInsertId()
						if err != nil {
							log.Fatal(err)
						}
						log.Printf("Inserted row with ID %d in DisLikes", id)
					}
					
				}
				cp.Put(ch)
			}
		}()
	}
	select{}
}