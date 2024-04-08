package storage

import (
	"slices"

	"github.com/emehrkay/cbt/pkg/types"
)

type Storage interface {
	Login(email, password string) (*types.User, error)
	GetUserByID(userID string) (*types.User, error)

	TicketCreate(ticket types.Ticket) (*types.Ticket, error)
	TicketDelete(ticket types.Ticket) error
	TicketUpdate(ticket types.Ticket) (*types.Ticket, error)
	TicketSearch(search TicketSearchParams) (types.TicketSet, error)
	TicketByID(ticketID string) (*types.Ticket, error)

	GetTrain() types.Train
	TrainDetails() (types.TicketSet, error)
	TrainCheckSeatTaken(car types.TrainCar, seat int32) (bool, error)
	TrainGetEmptySeats() (types.TicketSet, error)
}

type TicketSearchParams struct {
	IDs  []string
	Car  *types.TrainCar // needs to be {types.TrainCar, int}
	Seat *int32
	From []string
	To   []string
}

func (tsp TicketSearchParams) Match(all bool, ticket types.Ticket) bool {
	mismatches := 0
	clauses := 0

	if len(tsp.IDs) > 0 {
		clauses += 1

		if slices.Contains(tsp.IDs, ticket.ID) {
			mismatches += 1
		}
	}

	if tsp.Seat != nil && tsp.Car != nil {
		clauses += 1

		if ticket.Car == *tsp.Car && ticket.Seat == *tsp.Seat {
			mismatches += 1
		}
	}

	if len(tsp.From) > 0 {
		if all {
			clauses += len(tsp.From)
		} else {
			clauses += 1
		}

		if slices.Contains(tsp.From, ticket.From) {
			mismatches += 1
		}
	}

	if len(tsp.To) > 0 {
		if all {
			clauses += len(tsp.To)
		} else {
			clauses += 1
		}

		if slices.Contains(tsp.To, ticket.To) {
			mismatches += 1
		}
	}

	if all {
		return mismatches == clauses
	}

	return mismatches > 0
}
