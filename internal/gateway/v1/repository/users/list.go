package users

import (
	"context"
	"fmt"

	pb "github.com/guilherme-de-marchi/coin-commerce/api/users/v1"
	"github.com/guilherme-de-marchi/coin-commerce/pkg"
)

func List(ctx context.Context, req *pb.ListRequest) error {
	// body, err := proto.Marshal(req)
	// if err != nil {
	// 	return pkg.Error(err)
	// }

	// println("a: ", string(body))

	respChan := make(chan []byte, 1)
	pkg.Globals.MessageBroker.RequestChannel <- pkg.LoadBalancerRequest{
		Exchange:     "users",
		Target:       "List",
		Data:         req,
		ResponseChan: respChan,
	}

	fmt.Println(string(<-respChan))

	return nil
}
