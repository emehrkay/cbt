package rpc

import (
	"context"

	rpcBase "github.com/emehrkay/cbt/pkg/rpc"
)

func (s *server) Login(ctx context.Context, in *rpcBase.LoginRequest) (*rpcBase.LoginResponse, error) {

	login, err := s.service.Login(ctx, in.Email, in.Password)
	if err != nil {
		return nil, err
	}

	resp := rpcBase.LoginResponse{
		Jwt: login.JWT,
	}

	return &resp, nil
}
