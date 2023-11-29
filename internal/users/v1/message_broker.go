package v1

import (
	"context"
	"errors"
	"reflect"
	"sync"

	"github.com/guilherme-de-marchi/coin-commerce/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

type MessageBroker struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewMessageBroker() *MessageBroker {
	mb := &MessageBroker{}
	return mb
}

func (mb *MessageBroker) Run(s service) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return pkg.Error(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return pkg.Error(err)
	}

	err = ch.ExchangeDeclare(
		"users",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return pkg.Error(err)
	}

	mb.Conn = conn
	mb.Channel = ch

	return mb.setupQueues(s)
}

func (mb *MessageBroker) setupQueues(s service) error {
	wg := new(sync.WaitGroup)
	t := reflect.TypeOf(s)
	for i := 0; i < t.NumMethod(); i++ {
		ii := i
		q, err := mb.Channel.QueueDeclare(
			t.Method(ii).Name,
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return pkg.Error(err)
		}

		err = mb.Channel.QueueBind(
			q.Name,
			q.Name,
			"users",
			false,
			nil,
		)
		if err != nil {
			return pkg.Error(err)
		}

		wg.Add(1)
		go func() {
			if err := mb.consumeQueue(s, q.Name, t.Method(ii).Func); err != nil {
				log.Error().
					Any("queue", q.Name).
					Any("error", err).
					Msg("")
			}
			wg.Done()
		}()
	}

	wg.Wait()

	return nil
}

func (mb *MessageBroker) consumeQueue(s service, queueName string, f reflect.Value) error {
	deliveries, err := mb.Channel.Consume(
		queueName,
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return pkg.Error(err)
	}

	for d := range deliveries {
		dd := d
		go func() {
			if err := mb.handleDeliveries(s, dd, f); err != nil {
				log.Error().
					Any("queue", queueName).
					Any("error", err).
					Msg("")
			}
		}()
	}

	return nil
}

func (mb *MessageBroker) handleDeliveries(s service, d amqp.Delivery, f reflect.Value) error {
	if f.Kind() != reflect.Func {
		return pkg.Error(errors.New("reflect value is not a func"))
	}

	err := d.Ack(true)
	if err != nil {
		return pkg.Error(err)
	}

	resp := f.Call([]reflect.Value{
		reflect.ValueOf(s),
		reflect.ValueOf(d.Body),
	})

	if len(resp) != 2 {
		return pkg.Error(errors.New("unexpected behavior"))
	}

	bodyx, errx := resp[0], resp[1]
	if !errx.CanInterface() {
		return pkg.Error(errors.New("unexpected behavior"))
	}

	if errx.Interface() != nil {
		err, ok := errx.Interface().(error)
		if !ok {
			return pkg.Error(errors.New("unexpected behavior"))
		}

		return pkg.Error(err)
	}

	body := bodyx.Bytes()

	err = mb.Channel.PublishWithContext(
		context.Background(),
		"",
		d.ReplyTo,
		false,
		false,
		amqp.Publishing{
			ContentType:   "plain/text",
			Body:          body,
			CorrelationId: d.CorrelationId,
		},
	)
	if err != nil {
		return pkg.Error(err)
	}

	return nil
}
