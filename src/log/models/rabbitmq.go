package models

import (
  "log"
  "time"
  // "os"

  "github.com/streadway/amqp"
  "github.com/astaxie/beego"
  bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
  "encoding/json"


)

var driverNeo4j = bolt.NewDriver()

type TypeLogData struct {
	From string `json:"from"`
	To string `json:"to"`
	ReqId string `json:"reqId"`
	Header string `json:"header"`
	Body string `json:"body"`
	TypeRel string `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

func init() {

}

func CredMq() string {
	mq := ""
	if envOs == "local" {
		mq = beego.AppConfig.String("mq::local")
	}else if envOs == "dev" {
		mq = beego.AppConfig.String("mq::dev")
	}else if envOs == "prod" {
		mq = beego.AppConfig.String("mq::prod")
	}

	return mq
}

func CredNeo4j() string {
	neo4j := ""
	if envOs == "local" {
		neo4j = beego.AppConfig.String("neo4j::local")
	}else if envOs == "dev" {
		neo4j = beego.AppConfig.String("neo4j::dev")
	}else if envOs == "prod" {
		neo4j = beego.AppConfig.String("neo4j::prod")
	}

	return neo4j
}

func Receiver() {
	mq := CredMq()
	conn, err := amqp.Dial(mq)
	CheckErr(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	CheckErr(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"log", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	CheckErr(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	CheckErr(err, "Failed to register a consumer")

	forever := make(chan bool)
	var typeLogData TypeLogData

	go func() {
		for d := range msgs {
			err = json.Unmarshal(d.Body,&typeLogData);CheckErr(err,"error rabbitmq line 65")
			reqId := typeLogData.ReqId
			body := typeLogData.Body
			header := typeLogData.Header
			fromService := typeLogData.From+"_"+reqId
			toService := typeLogData.To+"_"+reqId
			nameFromService := typeLogData.From+"Service"
			nameToService := typeLogData.To+"Service"
			nmFromService := typeLogData.From+"Nm"
			nmToService := typeLogData.To+"Nm"



			
			if typeLogData.TypeRel == "req" {
				neoQuery := "create ("+nmFromService+":"+nameFromService+"{data:{data},reqId:{reqId},body:{body},header:{header},created_at:timestamp() })"
				execQuery := map[string]interface{}{"data":fromService,"reqId":reqId,"body":body,"header":header}
				SendGraphDb(neoQuery,execQuery)

				neoQuery2 := "create ("+nmToService+":"+nameToService+"{data:{data},reqId:{reqId},body:{body},header:{header},created_at:timestamp() })"
				execQuery2 := map[string]interface{}{"data":toService,"reqId":reqId,"body":"","header":""}
				SendGraphDb(neoQuery2,execQuery2)

				neoQuery3 := "match ("+nmFromService+":"+nameFromService+"{data:'"+fromService+"'}) "
				neoQuery3 += "match ("+nmToService+":"+nameToService+"{data:'"+toService+"'}) "
				neoQuery3 += "create ("+nmFromService+")-[:REQ]->("+nmToService+")"
				SendGraphPipeline(neoQuery3)
			}else if typeLogData.TypeRel == "res" {
				neoQuery := "merge ("+nmToService+":"+nameToService+"{data:{data}}) on match set "+nmToService+".body={body} "
				execQuery := map[string]interface{}{"data":toService,"body":body}
				SendGraphDb(neoQuery,execQuery)

				neoQuery3 := "match ("+nmFromService+":"+nameFromService+"{data:'"+fromService+"'}) "
				neoQuery3 += "match ("+nmToService+":"+nameToService+"{data:'"+toService+"'}) "
				neoQuery3 += "create ("+nmFromService+")-[:RES]->("+nmToService+")"
				SendGraphPipeline(neoQuery3)
			}

			SendDbLog(typeLogData)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}


func SendGraphDb(neoQuery string,execQuery map[string]interface{}) {
	credNeo4j := CredNeo4j()
	conn, err := driverNeo4j.OpenNeo(credNeo4j)
	CheckErr(err,"Error Connect GraphDB")
	defer conn.Close()

	stmt,err := conn.PrepareNeo(neoQuery)
	CheckErr(err,"error neo4j line 66")

	res,err := stmt.ExecNeo(execQuery)
	CheckErr(err,"error neo4j line 69")

	numRes,err := res.RowsAffected()
	CheckErr(err,"error neo4j line 72")
	beego.Debug("Created Rows")
	beego.Debug(numRes)

	stmt.Close()
}

func SendGraphPipeline(neoQuery string) {
	credNeo4j := CredNeo4j()
	conn, err := driverNeo4j.OpenNeo(credNeo4j)
	CheckErr(err,"Error Connect GraphDB")
	defer conn.Close()

	pipeline, err := conn.PreparePipeline(neoQuery)
	pipelineResults, err := pipeline.ExecPipeline(nil)
	CheckErr(err,"error neo4j line 125")
	for _, result := range pipelineResults {
		numResult, _ := result.RowsAffected()
		log.Println(numResult)
	}
}

type TypeLogDataDB struct {
	From string `json:"from"`
	To string `json:"to"`
	ReqId string `json:"reqId"`
	Header interface{} `json:"header"`
	Body interface{} `json:"body"`
	TypeRel string `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}
func SendDbLog(body TypeLogData) {
	var insHeader,insBody interface{}
	json.Unmarshal([]byte(body.Header),&insHeader)
	json.Unmarshal([]byte(body.Body),&insBody)

	insData := &TypeLogDataDB{
		From:body.From,
		To:body.To,
		ReqId:body.ReqId,
		Header:insHeader,
		Body:insBody,
		TypeRel:body.TypeRel,
		CreatedAt:time.Now(),
	}

	session := ConnectMongo();defer session.Close()

	c := session.DB("log").C("log")
	err := c.Insert(insData);CheckErr(err,"error logging mongoDB")

	return
}



