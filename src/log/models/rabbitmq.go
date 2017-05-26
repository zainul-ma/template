package models

import (
  "log"

  "github.com/streadway/amqp"
  "github.com/astaxie/beego"
  bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

func init() {

}

func Receiver() {
  conn, err := amqp.Dial("amqp://guest:guest@192.168.99.100:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"log", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
  <-forever
}

var driver = bolt.NewDriver()
func SendGraphDb(query string) {
  conn, err := driver.OpenNeo("bolt://localhost:7687")
	FailOnError(err,"Error Connect GraphDB")
	defer conn.Close()


  // stmt,err := conn.PrepareNeo("create")


}

func FailOnError(err error, msg string) {
	if err != nil {
		beego.Warning(msg)
    beego.Critical(err)
	}
}
