package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/emehrkay/cbt/pkg/types"
)

func (t *Train) BuyTicket(ctx context.Context, payloadd types.BuyTicket) (*types.Ticket, error) {
	auth, err := t.newAuth(ctx)
	if err != nil {
		return nil, err
	}

	train := t.store.GetTrain()

	// get the next open seat if either the car is empty. seat can be zero
	if strings.TrimSpace(string(payloadd.Car)) == "" {
		openSeats, err := t.TrainNextOpenSeat(ctx)
		if err != nil {
			return nil, err
		}

		if !openSeats.HaveOpenings() {
			return nil, types.ErrTrainFull
		}

		car, seat := openSeats.Pop()
		payloadd.Seat = *seat
		payloadd.Car = *car
	}

	ticket, err := t.store.TicketCreate(types.Ticket{
		UserID: auth.User.ID,
		Seat:   payloadd.Seat,
		Car:    payloadd.Car,
		To:     train.To,
		From:   train.From,
		Price:  train.TicketPrice,
	})
	if err != nil {
		return nil, fmt.Errorf(`unable to buy ticket -- %w`, err)
	}

	t.capture.Change(auth.User, "bought ticket", *ticket)

	return ticket, nil
}

func (t *Train) RemoveTicket(ctx context.Context, payload types.RemoveTicket) error {
	auth, err := t.newAuth(ctx)
	if err != nil {
		return err
	}

	ticket, err := t.store.TicketByID(payload.TicketID)
	if err != nil {
		return types.ErrNotFound
	}

	canEdit, err := auth.CanEditTicket(*ticket)
	if err != nil {
		return err
	}

	if !canEdit {
		return types.ErrInvalidAccess
	}

	err = t.store.TicketDelete(*ticket)
	if err == nil {
		t.capture.Change(auth.User, "deleted ticket", ticket)
	}

	return err
}

func (t *Train) ChangeTicket(ctx context.Context, payload types.ChangeTicket) (*types.Ticket, error) {
	auth, err := t.newAuth(ctx)
	if err != nil {
		return nil, err
	}

	ticket, err := t.TicketByID(ctx, payload.TicketID)
	if err != nil {
		return nil, types.ErrNotFound
	}

	canEdit, err := auth.CanEditTicket(*ticket)
	if err != nil {
		return nil, err
	}

	if !canEdit {
		return nil, types.ErrInvalidAccess
	}

	updatedTicket := *ticket // dont change item ref'd in memory
	updatedTicket.Car = payload.Car
	updatedTicket.Seat = payload.Seat
	updated, err := t.store.TicketUpdate(updatedTicket)
	if err != nil {
		return nil, fmt.Errorf(`unable to update ticket -- %w`, err)
	}

	t.capture.Change(auth.User, "udpated ticket", updated)

	return updated, nil
}

func (t *Train) TicketByID(ctx context.Context, ticketID string) (*types.Ticket, error) {
	_, err := t.newAuth(ctx)
	if err != nil {
		return nil, err
	}

	return t.store.TicketByID(ticketID)
}
