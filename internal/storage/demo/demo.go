package demo

import (
	"fmt"

	"github.com/emehrkay/cbt/internal/storage"
	"github.com/emehrkay/cbt/pkg/types"
)

type DemoStorage struct{}

func (ds DemoStorage) GetUserByID(userID string) (*types.User, error) {
	for _, user := range demoUsers {
		if user.ID == userID {
			return user, nil
		}
	}

	return nil, types.ErrNotFound
}

func (ds DemoStorage) Login(email, password string) (*types.User, error) {
	for _, user := range demoUsers {
		if user.Email == email && user.Password == password {
			return user, nil
		}
	}

	return nil, types.ErrNotFound
}

func (ds DemoStorage) TicketCreate(ticket types.Ticket) (*types.Ticket, error) {
	// in a proper db, this would be taken care of with constraints
	if taken, _ := ds.TrainCheckSeatTaken(ticket.Car, ticket.Seat); taken {
		return nil, types.ErrSeatTaken
	}

	tickets, ok := demoTrain[ticket.Car]
	if !ok {
		return nil, types.ErrInvalidCar
	}

	if len(tickets) >= maxSeatCount {
		return nil, types.ErrInvalidSeat
	}

	ticket.ID = fmt.Sprintf(`%v`, len(demoTrain[ticket.Car]))
	demoTrain[ticket.Car] = append(demoTrain[ticket.Car], &ticket)

	return &ticket, nil
}

func (ds DemoStorage) TicketDelete(ticket types.Ticket) error {
	tickets, ok := demoTrain[ticket.Car]
	if !ok {
		return types.ErrNotFound
	}

	for i, tic := range demoTrain[ticket.Car] {
		if ticket.ID == tic.ID {
			tickets = append(tickets[:i], tickets[i+1:]...)
		}
	}

	demoTrain[ticket.Car] = tickets

	return nil
}

func (ds DemoStorage) TicketUpdate(updatedTicket types.Ticket) (*types.Ticket, error) {
	ticket, err := ds.TicketByID(updatedTicket.ID)
	if err != nil {
		return nil, err
	}

	tickets, ok := demoTrain[ticket.Car]
	if !ok {
		return nil, types.ErrInvalidCar
	}

	for i, tic := range tickets {
		if tic != nil && tic.ID == ticket.ID {
			err := ds.TicketDelete(*ticket)
			if err != nil {
				return nil, err
			}

			tickets[i] = &updatedTicket
			demoTrain[updatedTicket.Car] = tickets
			return &updatedTicket, nil
		}
	}

	return nil, types.ErrNotFound
}

func (ds DemoStorage) TicketSearch(search storage.TicketSearchParams) (types.TicketSet, error) {
	var found types.TicketSet

	for _, tickets := range demoTrain {
		for _, ticket := range tickets {
			add := search.Match(true, *ticket)
			if add {
				found = append(found, ticket)
			}
		}
	}

	if len(found) == 0 {
		return nil, types.ErrNotFound
	}

	return found, nil
}

func (ds DemoStorage) TicketByID(ticketID string) (*types.Ticket, error) {
	tickets, err := ds.TicketSearch(storage.TicketSearchParams{
		IDs: []string{ticketID},
	})
	if err != nil {
		return nil, err
	}

	return tickets[0], nil
}

func (ds DemoStorage) TrainDetails() (types.TicketSet, error) {
	var tickets types.TicketSet

	for _, tics := range demoTrain {
		tickets = append(tickets, tics...)
	}

	return tickets, nil
}

func (ds DemoStorage) TrainCheckSeatTaken(car types.TrainCar, seat int32) (bool, error) {
	_, ok := demoTrain[car]
	if !ok {
		return false, types.ErrInvalidCar
	}

	search := storage.TicketSearchParams{
		Seat: &seat,
		Car:  &car,
	}
	found, err := ds.TicketSearch(search)
	if err != nil {
		return false, err
	}

	return found != nil, nil
}

func (ds DemoStorage) TrainGetEmptySeats() (types.TicketSet, error) {
	taken := map[types.TrainCar][10]bool{}

	for car, tickets := range demoTrain {
		t := taken[car]
		for _, ticket := range tickets {
			if ticket == nil {
				continue
			}

			t[ticket.Seat] = true
		}

		taken[car] = t
	}

	open := types.TicketSet{}

	for car, seated := range taken {
		for i, t := range seated {
			if !t {
				open = append(open, &types.Ticket{
					Car:  car,
					Seat: int32(i),
				})
			}
		}
	}

	return open, nil
}

func (ds DemoStorage) GetTrain() types.Train {
	return train
}
