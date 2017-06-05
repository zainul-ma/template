package thirdparty

import(
	"github.com/streadway/amqp"
 	"github.com/astaxie/beego"
	"log"
	"encoding/json"
)
type TypeLogData struct {
	From string `json:"from"`
	To string `json:"to"`
	ReqId string `json:"reqId"`
	Header string `json:"header"`
	Body string `json:"body"`
	TypeRel string `json:"type"`
}

func SendMq(inputReqBody []byte,fromService string,toService string,headerAll string,reqId string,newRequest bool,typeRelation string) {
 	log.Println("")

 	conn, err := amqp.Dial("amqp://guest:guest@172.17.0.1:5672/")
	CheckErr(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	CheckErr(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"log",	// name
		false,	// durable
		false,	// delete when unused
		false,	// exclusive
		false,	// no-wait
		nil,	// arguments
	)
	CheckErr(err, "Failed to declare a queue")

	typeLogData := &TypeLogData{
		From:fromService,
		To:toService,
		ReqId:reqId,
		Header:headerAll,
		Body:string(inputReqBody),
		TypeRel:typeRelation,
	}
	logJsonMarshal,err := json.Marshal(typeLogData);CheckErr(err,"error rabbitMQ line 47")
	
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        logJsonMarshal,
		},
	)
	CheckErr(err, "Failed to publish a message")
}

func CheckErr(err error, msg string) {
	if err != nil {
	    beego.Warning("MESSAGE RABBITMQ: "+msg)
	    beego.Error(1024,err)
	}
}
