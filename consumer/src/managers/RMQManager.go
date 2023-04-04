package managers

import(
	"log"
	"goTwinderRMQConsumer/src/tools"
	"goTwinderRMQConsumer/src/helpers"
	amqp "github.com/rabbitmq/amqp091-go"
)

func RMQConsumeWithQName( channelPool *tools.ChannelPool, qName string, threadNum int) <-chan amqp.Delivery {
	ch := channelPool.Get()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C thread %d", threadNum)
	q, err := ch.QueueDeclare(
		"swipeQueue", // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
		)
	helpers.FailOnError(err, "Failed to declare a queue", "ch.Consume()")
	
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
		)
	helpers.FailOnError(err, "Failed to register a consumer", "ch.Consume()")
	return msgs
}