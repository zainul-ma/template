package thirdparty

import (
	"encoding/json"
	"os"

	"github.com/astaxie/beego"
	"github.com/streadway/amqp"
)

// TypeLogData = request Log Data
type TypeLogData struct {
	From    string `json:"from"`
	To      string `json:"to"`
	ReqID   string `json:"reqId"`
	Header  string `json:"header"`
	Body    string `json:"body"`
	TypeRel string `json:"type"`
}

// CredMq = Get Credential MQ
func CredMq() string {
	mq := ""
	envOs := os.Getenv("GOENV")
	if envOs == "local" {
		mq = beego.AppConfig.String("mq::local")
	} else if envOs == "dev" {
		mq = beego.AppConfig.String("mq::dev")
	} else if envOs == "prod" {
		mq = beego.AppConfig.String("mq::prod")
	}

	return mq
}

func ConnectMq(mq string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(mq)
	CheckErr(err, "Failed to connect to RabbitMQ")

	return conn, err
}

// SendMq = to Send MQ Data
func SendMq(inputReqBody []byte, fromService string, toService string, headerAll string, reqID string, newRequest bool, typeRelation string) {
	mq := CredMq()

	conn, err := ConnectMq(mq)
	CheckErr(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	CheckErr(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"log", // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	CheckErr(err, "Failed to declare a queue")

	typeLogData := &TypeLogData{
		From:    fromService,
		To:      toService,
		ReqID:   reqID,
		Header:  headerAll,
		Body:    string(inputReqBody),
		TypeRel: typeRelation,
	}
	logJSONMarshal, err := json.Marshal(typeLogData)
	CheckErr(err, "error rabbitMQ line 47")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        logJSONMarshal,
		},
	)
	CheckErr(err, "Failed to publish a message")

	return
}

// CheckErr = Checking Error
func CheckErr(err error, msg string) {
	if err != nil {
		beego.Warning("MESSAGE RABBITMQ: " + msg)
		beego.Error(1024, err)
	}
}
