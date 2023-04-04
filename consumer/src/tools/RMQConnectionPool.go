package tools

import (
	"sync"
	"log"
	amqp "github.com/rabbitmq/amqp091-go"
)

// A ChannelPool represents a pool of channels.
type ChannelPool struct {
	connection *amqp.Connection
	mu     sync.Mutex
	chPool chan *amqp.Channel
}

// NewChannelPool creates a new ChannelPool with a given capacity and a factory function that produces new items for the pool.
func NewRMQChannelPool(capacity int) *ChannelPool {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	chPool := make(chan *amqp.Channel, capacity)
	for i:= 0; i < capacity; i++ {
		ch, err := conn.Channel()
		if (err != nil) {
			failOnError(err, "Error initializing channels")
			continue
		}
		chPool <- ch
	}

	return &ChannelPool{
		connection: conn,
		chPool: chPool,
	}
}

// Get returns a channel from the pool.
func (cp *ChannelPool) Get() *amqp.Channel {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	select {
	case ch := <-cp.chPool:
		return ch
	default:
		ch, err := cp.connection.Channel()
		if (err != nil) {
			failOnError(err, "Error re-creating channel")
			return nil
		}
		return ch
	}
}

// Put returns a channel to the pool.
func (cp *ChannelPool) Put(ch *amqp.Channel) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.chPool <- ch
}


func failOnError(err error, msg string) {
	if err != nil {
	  log.Panicf("%s: %s", msg, err)
	}
}