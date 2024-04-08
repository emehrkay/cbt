package rpc

import (
	"context"

	trainCTX "github.com/emehrkay/cbt/internal/models/jwt/ctx"
	rpcBase "github.com/emehrkay/cbt/pkg/rpc"
)

func (s *server) TrainOpenSeats(ctx context.Context, in *rpcBase.TrainOpenSeatRequest) (*rpcBase.TrainOpenSeatResponses, error) {
	token, err := s.tokenFromContextMetadata(ctx)
	if err != nil {
		return nil, err
	}

	ctx = trainCTX.ContextWithToken(ctx, *token)
	open, err := s.service.TrainNextOpenSeat(ctx)
	if err != nil {
		return nil, err
	}

	resp := &rpcBase.TrainOpenSeatResponses{
		Seats: []*rpcBase.TrainOpenSeatResponse{},
	}

	for car, seats := range open.Seats {
		resp.Seats = append(resp.Seats, &rpcBase.TrainOpenSeatResponse{
			Car:   string(car),
			Seats: seats,
		})
	}

	return resp, nil
}

func (s *server) TrainDetails(ctx context.Context, in *rpcBase.TrainDetailsRequest) (*rpcBase.TrainDetailsResponse, error) {
	token, err := s.tokenFromContextMetadata(ctx)
	if err != nil {
		return nil, err
	}

	ctx = trainCTX.ContextWithToken(ctx, *token)
	details, err := s.service.TrainDetails(ctx)
	if err != nil {
		return nil, err
	}

	resp := details.ToProto()

	return resp, nil
}
