package users

import (
	"context"

	pb "github.com/guilherme-de-marchi/coin-commerce/api/users/v1"
	"github.com/guilherme-de-marchi/coin-commerce/pkg"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

func List(ctx context.Context, req *pb.ListRequest) error {
	ch, err := pkg.Globals.RabbitMQ.Channel()
	if err != nil {
		return pkg.Error(err)
	}

	body, err := proto.Marshal(&pb.ListRequest{})
	if err != nil {
		return pkg.Error(err)
	}

	err = ch.PublishWithContext(
		ctx,
		"",
		"users",
		false,
		false,
		amqp091.Publishing{
			ContentType:   "plain/text",
			Body:          body,
			ReplyTo:       "gateway",
			CorrelationId: "123456",
		},
	)
	if err != nil {
		return pkg.Error(err)
	}

	return nil
}
