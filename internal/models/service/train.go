package service

import (
	"context"

	"github.com/emehrkay/cbt/pkg/types"
)

func (t *Train) TrainDetails(ctx context.Context) (*types.TrainDetails, error) {
	auth, err := t.newAuth(ctx)
	if err != nil {
		return nil, err
	}

	if !auth.IsAdmin() {
		return nil, types.ErrInvalidAccessAdminOnly
	}

	tickets, err := t.store.TrainDetails()
	if err != nil {
		return nil, err
	}

	resp := types.TrainDetails{}

	for _, ticket := range tickets {
		user, err := t.store.GetUserByID(ticket.UserID)
		if err != nil {
			return nil, err
		}

		resp.Trains = append(resp.Trains, types.TrainDetail{
			Train:  t.store.GetTrain(),
			User:   *user,
			Ticket: *ticket,
		})
	}

	return &resp, nil
}

func (t *Train) TrainNextOpenSeat(ctx context.Context) (*types.TrainOpenSeatsResp, error) {
	_, err := t.newAuth(ctx)
	if err != nil {
		return nil, err
	}

	open, err := t.store.TrainGetEmptySeats()
	if err != nil {
		return nil, err
	}

	resp := types.NewTrainOpenSeatsResp()
	resp.AddTicket(open...)

	return resp, nil
}
