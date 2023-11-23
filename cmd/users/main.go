package main

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"

	users "github.com/guilherme-de-marchi/coin-commerce/api/users/v1"
)

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

	c, err := ch.QueueDeclare("carlinhos", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msg, err := proto.Marshal(&users.CreateRequest{
		Name:  "carlos",
		Email: "maria@gmail",
		Phone: "abc123",
	})
	if err != nil {
		panic(err)
	}

	err = ch.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType: "plain/text",
		Body:        msg,
		ReplyTo:     "carlinhos",
	})
	if err != nil {
		panic(err)
	}

	deliveries, err := ch.ConsumeWithContext(ctx, c.Name, "", true, false, false, false, nil)
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
	}

	println("sent")
}
