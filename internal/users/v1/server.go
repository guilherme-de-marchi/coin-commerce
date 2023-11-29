package v1

import (
	"context"
	"fmt"

	pb "github.com/guilherme-de-marchi/coin-commerce/api/users/v1"
	"github.com/guilherme-de-marchi/coin-commerce/pkg"
	"google.golang.org/protobuf/proto"
)

type service interface {
	List([]byte) ([]byte, error)
	Get([]byte) ([]byte, error)
	Create([]byte) ([]byte, error)
	Update([]byte) ([]byte, error)
	Delete([]byte) ([]byte, error)
}

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) List(rawIn []byte) ([]byte, error) {
	var req pb.ListRequest
	err := proto.Unmarshal(rawIn, &req)
	if err != nil {
		return nil, pkg.Error(err)
	}

	out, err := s.list(context.Background(), &req)
	if err != nil {
		return nil, pkg.Error(err)
	}

	rawOut, err := proto.Marshal(out)
	if err != nil {
		return nil, pkg.Error(err)
	}

	return rawOut, nil
}

func (s *Server) list(ctx context.Context, in *pb.ListRequest) (*pb.ListResponse, error) {
	fmt.Println(in)
	return nil, nil
}

func (s *Server) Get(rawIn []byte) ([]byte, error) {
	return nil, nil
}

func (s *Server) Create(rawIn []byte) ([]byte, error) {
	return nil, nil
}

func (s *Server) Update(rawIn []byte) ([]byte, error) {
	return nil, nil
}

func (s *Server) Delete(rawIn []byte) ([]byte, error) {
	return nil, nil
}
