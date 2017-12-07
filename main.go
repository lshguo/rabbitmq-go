package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
	"strconv"
	"sync"
)

const rabbitUrl  = "amqp://guest:guest@192.168.241.130:5672/"
const exType     = "topic"
const exName     = "excellent"

var number uint32
var numberLock 	sync.Mutex

func errpanic(err error){
	if err != nil{
		panic(err)
	}
}
func applyNumber()uint32{
	numberLock.Lock()

	ret := number
	number = number + 1

	numberLock.Unlock()

	return ret
}
func createExchange(rabbitUrl, exName, exType string){
	conn, err := amqp.Dial(rabbitUrl)
	errpanic(err)
	defer conn.Close()

	ch, err := conn.Channel()
	errpanic(err)
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exName,
		exType,
		true,
		false,
		false,
		false,
		nil)
	errpanic(err)

	fmt.Println("create ex: ", exName)
}

func send(rabbitUrl, exName string){
	conn, err := amqp.Dial(rabbitUrl)
	errpanic(err)
	defer conn.Close()

	ch, err := conn.Channel()
	errpanic(err)
	defer ch.Close()

	for i := 0;;i++{
		key := ([]string{"co.black", "co.red","log.black","log.red"})[i%4]
		body := strconv.Itoa(i)
		body = body + " " + key
		err = ch.Publish(
			exName,
			key,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body: []byte(body),
			})
		errpanic(err)
		fmt.Println("send msg ", body)
		time.Sleep(time.Second*2)
	}
}

func recv(rabbitUrl, exName string){
	index := applyNumber()
	routingKey := ([]string{"co.*","*.black","log.*","*.red"})[index%4]
	//consumerName := strconv.Itoa(number)

	conn, err := amqp.Dial(rabbitUrl)
	errpanic(err)
	defer conn.Close()

	ch, err := conn.Channel()
	errpanic(err)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",
		/* durable和autoDelete的值最好相反，这样特性比较一致*/
		false,
		/* autoDelete=true: queue will be del when conn closed */
		true,
		/* exclusive=true: queue will be del when conn closed */
		false,
		false,
		nil)
	errpanic(err)

	err = ch.QueueBind(
		q.Name,
		routingKey,
		 exName,
		 false,
		 nil)
	errpanic(err)

	msgs, err := ch.Consume(
		q.Name,
		routingKey,
		true,
		true,
		false,
		false,
		nil)
	errpanic(err)

	for {
		if msg, ok := <- msgs;ok{
			fmt.Println("    ", routingKey, " recv msg ", string(msg.Body))
		}else {
			fmt.Println("channel ", routingKey," closed")
			break
		}
	}
}
func main(){
	createExchange(rabbitUrl, exName, exType)
	go send(rabbitUrl, exName)
	go recv(rabbitUrl, exName)
	go recv(rabbitUrl, exName)
	select {
	//
	}
}