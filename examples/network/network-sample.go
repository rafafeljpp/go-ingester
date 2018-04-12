// main package.
// Solo para probar el perfomance
package main

import (
	"fmt"
	"go-ingester/pkg/pool"
	"log"
	"runtime"
	"time"

	"github.com/firstrow/tcp_server"
	"github.com/streadway/amqp"
)

// Job mi trabajo
type Job struct {
	payload string
	init    time.Time
	ip      string
	port    int
	delay   time.Duration
}

var cantidad int
var conn *amqp.Connection
var ch *amqp.Channel
var q amqp.Queue
var flag bool

func main() {
	runtime.GOMAXPROCS(8)

	p := pool.NewManager(10, 5)
	server := tcp_server.New("127.0.0.1:9999")

	p.Start()
	defer p.Stop()
	defer conn.Close()

	server.OnNewClient(func(c *tcp_server.Client) {

		//fmt.Printf("Nueva conexión desde: %s", c.Conn().RemoteAddr())
	})
	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		job := Job{message, time.Now(), "", 0, 0}
		p.AddJob(job)
	})

	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		// connection with client lost
	})

	server.Listen()

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// IsValid Método para validar el Job
func (j Job) IsValid() bool {

	return true
}

// Publish  donde yo quiera
func (j Job) Publish() (bool, error) {

	var err error
	if !flag {
		conn, err = amqp.Dial("amqp://guest:guest@rbm-local.sigis.com.ve:5672/")
		failOnError(err, "Failed to connect to RabbitMQ")
		ch, err = conn.Channel()
		failOnError(err, "Failed to declare a queue")
		flag = true
	}

	q, err = ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		true,    // no-wait
		nil,     // arguments
	)
	//failOnError(err, "Failed to declare a queue")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(j.payload),
		})
	//failOnError(err, "Failed to publish a message")

	return true, nil
}

// Rejected Recibir los mensajes rechazados
func (j Job) Rejected(e error) {
	var elapse time.Duration
	elapse = time.Since(j.init)

	fmt.Println("Mensaje Rechazado... ", elapse)
}
