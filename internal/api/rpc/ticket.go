package rpc

import (
	"context"

	trainCTX "github.com/emehrkay/cbt/internal/models/jwt/ctx"
	rpcBase "github.com/emehrkay/cbt/pkg/rpc"
	"github.com/emehrkay/cbt/pkg/types"
)

func (s *server) BuyTicket(ctx context.Context, in *rpcBase.BuyTicketRequest) (*rpcBase.TicketResponse, error) {
	token, err := s.tokenFromContextMetadata(ctx)
	if err != nil {
		return nil, err
	}

	ctx = trainCTX.ContextWithToken(ctx, *token)
	ticket, err := s.service.BuyTicket(ctx, types.BuyTicket{})
	if err != nil {
		return nil, err
	}

	return ticket.ToProto(), nil
}

func (s *server) RemoveTicket(ctx context.Context, in *rpcBase.RemoveTicketRequest) (*rpcBase.Empty, error) {
	token, err := s.tokenFromContextMetadata(ctx)
	if err != nil {
		return nil, err
	}

	ctx = trainCTX.ContextWithToken(ctx, *token)
	err = s.service.RemoveTicket(ctx, types.RemoveTicket{
		TicketID: in.TicketId,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *server) ChangeTicket(ctx context.Context, in *rpcBase.ChangeTicketRequest) (*rpcBase.TicketResponse, error) {
	token, err := s.tokenFromContextMetadata(ctx)
	if err != nil {
		return nil, err
	}

	payload := &types.ChangeTicket{}
	payload.FromProto(in)

	ctx = trainCTX.ContextWithToken(ctx, *token)
	updated, err := s.service.ChangeTicket(ctx, *payload)
	if err != nil {
		return nil, err
	}

	return updated.ToProto(), nil
}
