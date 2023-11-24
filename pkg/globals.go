package pkg

import (
	"fmt"
	"reflect"

	amqp "github.com/rabbitmq/amqp091-go"
)

var Globals globals

type globals struct {
	RabbitMQ *amqp.Connection
}

func (g globals) Validate() error {
	t := reflect.ValueOf(g)
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).IsNil() {
			return Error(fmt.Errorf("field %s is nil", t.Type().Field(i).Name))
		}
	}
	return nil
}
