package thirdparty

import(
  "github.com/streadway/amqp"
  "github.com/astaxie/beego"
	"log"
)

func SendMQ(inputReqBody []byte,fromService string,toService string,headerAll string,reqID string) {
  conn, err := amqp.Dial("amqp://guest:guest@192.168.99.100:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"log", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	sendBuffer := `{
    "from":`+fromService+`,
    "to":`+toService+`,
    "reqID":`+reqID+`,
    "header":`+headerAll+`,
    "body":`+string(inputReqBody)+
  `}`
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(sendBuffer),
		})
	log.Printf(" [x] Sent %s", sendBuffer)
	failOnError(err, "Failed to publish a message")
}

func failOnError(err error, msg string) {
  if err != nil {
    beego.Warning("MESSAGE RABBITMQ: "+msg)
    beego.Error(1024,err)
	}
}
