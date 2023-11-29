package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type MessageBroker struct {
	Conn           *amqp091.Connection
	Prefix         string
	RequestChannel LoadBalancerRequestChannel
}

func NewMessageBroker(prefix string) (*MessageBroker, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, Error(err)
	}

	return &MessageBroker{
		Conn:           conn,
		Prefix:         prefix,
		RequestChannel: make(LoadBalancerRequestChannel, 1),
	}, nil
}

func (m *MessageBroker) Start(ctx context.Context, workers int) {
	lb := NewLoadBalancer(m.RequestChannel)
	lb.Start(ctx, workers, m.workerFunc)
}

func (m *MessageBroker) HandleWorkerFunc(ctx context.Context, wg *sync.WaitGroup, workerChan LoadBalancerRequestChannel) {
	if err := m.workerFunc(ctx, wg, workerChan); err != nil {
		log.Error().
			Any("message_broker", m).
			Any("error", err).
			Msg("")
	}
}

func (m *MessageBroker) workerFunc(ctx context.Context, wg *sync.WaitGroup, workerChan LoadBalancerRequestChannel) error {
	defer wg.Done()
	defer close(workerChan)

	ch, err := m.Conn.Channel()
	if err != nil {
		return Error(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		generateQueueName(m.Prefix),
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return Error(err)
	}

	correlations := NewAsyncMap[string, chan []byte]()

	wg2 := new(sync.WaitGroup)
	wg2.Add(2)

	go func() {
		if err := m.handleRequests(ctx, wg2, ch, workerChan, correlations, q.Name); err != nil {
			log.Error().
				Any("message_broker", m).
				Any("error", err).
				Msg("")
		}
	}()

	go func() {
		if err := m.handleDeliveries(ch, wg2, correlations, q.Name); err != nil {
			log.Error().
				Any("message_broker", m).
				Any("error", err).
				Msg("")
		}
	}()

	wg2.Wait()

	return nil
}

func (m *MessageBroker) handleRequests(ctx context.Context, wg *sync.WaitGroup, ch *amqp.Channel, workerChan LoadBalancerRequestChannel, correlations *AsyncMap[string, chan []byte], queueName string) error {
	defer wg.Done()

	for req := range workerChan {
		var body []byte
		var err error

		switch v := req.Data.(type) {
		case protoreflect.ProtoMessage:
			body, err = proto.Marshal(v)
		default:
			body, err = json.Marshal(v)
		}

		if err != nil {
			return Error(err)
		}

		correlation := uuid.NewString()

		err = ch.PublishWithContext(
			ctx,
			req.Exchange,
			req.Target,
			false,
			false,
			amqp091.Publishing{
				ContentType:   "plain/text",
				Body:          body,
				ReplyTo:       queueName,
				CorrelationId: correlation,
			},
		)
		if err != nil {
			return Error(err)
		}

		correlations.Set(correlation, req.ResponseChan)
	}

	return nil
}

func (m *MessageBroker) handleDeliveries(ch *amqp.Channel, wg *sync.WaitGroup, correlations *AsyncMap[string, chan []byte], queueName string) error {
	defer wg.Done()

	ds, err := ch.Consume(
		queueName,
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return Error(err)
	}

	for d := range ds {
		responseChan, ok := correlations.Get(d.CorrelationId)
		if !ok {
			continue
		}

		responseChan <- d.Body
		if err := d.Ack(true); err != nil {
			return Error(err)
		}
	}

	return nil
}

func generateQueueName(prefix string) string {
	return fmt.Sprintf("%s:%s", prefix, uuid.NewString())
}
