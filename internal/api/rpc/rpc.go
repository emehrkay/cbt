package rpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/emehrkay/cbt/internal/models/service"
	rpcBase "github.com/emehrkay/cbt/pkg/rpc"
)

func New(trainService *service.Train) (*server, error) {
	return &server{
		service: trainService,
	}, nil
}

type server struct {
	rpcBase.UnimplementedTrainServer
	service   *service.Train
	clentConn *grpc.ClientConn
}

func (s *server) Run(target string, options ...grpc.DialOption) error {
	var err error

	s.clentConn, err = grpc.NewClient(target, options...)
	if err != nil {
		return fmt.Errorf(`unable to create grpc connection to: %v -- %w`, target, err)
	}

	defer s.clentConn.Close()

	listener, err := net.Listen("tcp", target)
	if err != nil {
		return fmt.Errorf(`unable to create tpc conn: %v -- %w`, target, err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(s.service.GetValidator().UnaryServerInterceptor),
	)
	rpcBase.RegisterTrainServer(server, s)

	fmt.Printf("server listening at %v\n", listener.Addr())

	go func() {
		if err := server.Serve(listener); err != nil {
			msg := fmt.Sprintf(`unable to start grpc server -- %v`, err)
			panic(msg)
		}
	}()

	return nil
}

func (s *server) tokenFromContextMetadata(ctx context.Context) (*string, error) {
	// rip the token from the metadata via the context
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("no metadata found in context")
	}
	tokens := headers.Get("jwt")
	if len(tokens) < 1 {
		return nil, errors.New("no token found in metadata")
	}
	tokenString := tokens[0]

	return &tokenString, nil
}
