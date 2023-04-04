package main

import (
	"log"
	"goTwinderRMQConsumer/src/tools"
	// amqp "github.com/rabbitmq/amqp091-go"
)

var numChannel = 20
  
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
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
				}
				cp.Put(ch)
			}
		}()
	}
	select{}
}