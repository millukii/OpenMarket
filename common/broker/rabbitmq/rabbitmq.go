package rabbitmq

import (
	"fmt"
	"log"

	"github.com/millukii/commons/broker"
	amqp "github.com/rabbitmq/amqp091-go"
)
func Connect(user, pass, host, port string) (*amqp.Channel, func() error ){ 

	address := fmt.Sprintf("amqp://%s:%s@%s:%s", user, pass, host,port)

	con, err := amqp.Dial(address)

	if err != nil {
		log.Fatal(err)
	}

	ch, err := con.Channel()

	if err != nil {
		log.Fatal(err)
	}

	err = ch.ExchangeDeclare(broker.OrderCreatedEvent, "direct", true, false, false,false,nil)
	if err != nil {
		log.Fatal(err)
	}
	err = ch.ExchangeDeclare(broker.OrderCreatedPaid,"fanout", true,false,false,false,nil)
	if err != nil {
		log.Fatal(err)
	}
	return ch, con.Close
}