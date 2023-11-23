package main

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"

	users "github.com/guilherme-de-marchi/coin-commerce/api/users/v1"
)

type UserService struct {
	users.UnimplementedUserServiceServer
}

func (UserService) Create(context.Context, *users.CreateRequest) (*users.User, error) {
	return nil, nil
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("test", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	deliveries, err := ch.ConsumeWithContext(ctx, q.Name, "consumer 1", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	for d := range deliveries {
		msg := &users.CreateRequest{}
		err = proto.Unmarshal(d.Body, msg)
		if err != nil {
			panic(err)
		}

		fmt.Println(msg)

		err = ch.PublishWithContext(ctx, "", d.ReplyTo, false, false, amqp.Publishing{
			ContentType: "plain/text",
			Body:        d.Body,
		})
		if err != nil {
			panic(err)
		}
	}
}
