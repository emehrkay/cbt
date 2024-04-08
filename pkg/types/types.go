package types

import (
	"fmt"

	rpcBase "github.com/emehrkay/cbt/pkg/rpc"
)

type UserLevel string
type TrainCar string

var (
	LevelUser  UserLevel = "user"
	LevelAdmin UserLevel = "admin"
	CarA       TrainCar  = "a"
	CarB       TrainCar  = "b"
)

type UserSet []*User

type User struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Level     UserLevel `json:"level"`
}

func (u User) ToProto() *rpcBase.User {
	return &rpcBase.User{
		Id:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Level:     string(u.Level),
	}
}

func (u User) String() string {
	return fmt.Sprintf(`user[%s] %s %s`, u.ID, u.FirstName, u.LastName)
}

type TicketSet []*Ticket

type Ticket struct {
	ID     string   `json:"id"`
	UserID string   `json:"user_id"`
	From   string   `json:"from"`
	To     string   `json:"to"`
	Price  Currency `json:"price"`
	Car    TrainCar `json:"car"`
	Seat   int32    `json:"seat"`
}

func (t Ticket) ToProto() *rpcBase.TicketResponse {
	return &rpcBase.TicketResponse{
		Id:          t.ID,
		UserId:      t.UserID,
		Source:      t.From,
		Destination: t.To,
		Price:       float32(t.Price),
		Car:         string(t.Car),
		Seat:        t.Seat,
	}
}

type LoginResponse struct {
	JWT string `json:"jwt"`
}

type Train struct {
	Name        string   `json:"name"`
	From        string   `json:"from"`
	To          string   `json:"to"`
	TicketPrice Currency `json:"price"`
}

func (t Train) ToProto() *rpcBase.TrainResponse {
	return &rpcBase.TrainResponse{
		Name: t.Name,
		From: t.From,
		To:   t.To,
	}
}

func NewTrainOpenSeatsResp() *TrainOpenSeatsResp {
	return &TrainOpenSeatsResp{
		Seats: map[TrainCar][]int32{},
	}
}

type TrainOpenSeatsResp struct {
	Seats map[TrainCar][]int32 `json:"seats"`
}

func (t *TrainOpenSeatsResp) AddTicket(tickets ...*Ticket) {
	for _, ticket := range tickets {
		if ticket == nil {
			continue
		}

		_, ok := t.Seats[ticket.Car]
		if !ok {
			t.Seats[ticket.Car] = []int32{}
		}

		t.Seats[ticket.Car] = append(t.Seats[ticket.Car], ticket.Seat)
	}
}

func (t TrainOpenSeatsResp) HaveOpenings() bool {
	total := 0

	for _, seats := range t.Seats {
		total += len(seats)
	}

	return total > 0
}

func (t TrainOpenSeatsResp) Pop() (*TrainCar, *int32) {
	for car, seats := range t.Seats {
		return &car, &seats[0]
	}

	return nil, nil
}

type BuyTicket struct {
	Car  TrainCar `json:"car"`
	Seat int32    `json:"seat"`
}

type RemoveTicket struct {
	TicketID string `json:"ticket_id"`
}

type TrainDetail struct {
	Train  Train  `json:"train"`
	User   User   `json:"user"`
	Ticket Ticket `json:"ticket"`
}

func (t TrainDetail) ToProto() *rpcBase.TrainDetails {
	return &rpcBase.TrainDetails{
		Ticket: t.Ticket.ToProto(),
		Train:  t.Train.ToProto(),
		User:   t.User.ToProto(),
	}
}

type TrainDetails struct {
	Trains []TrainDetail `json:"trains"`
}

func (td TrainDetails) ToProto() *rpcBase.TrainDetailsResponse {
	proto := rpcBase.TrainDetailsResponse{}

	for _, train := range td.Trains {
		proto.Trains = append(proto.Trains, train.ToProto())
	}

	return &proto
}

type ChangeTicket struct {
	TicketID string   `json:"ticket_id"`
	Car      TrainCar `json:"car"`
	Seat     int32    `json:"seat"`
}

func (ct *ChangeTicket) FromProto(proto *rpcBase.ChangeTicketRequest) {
	ct.TicketID = proto.TicketId
	ct.Car = TrainCar(proto.Car)
	ct.Seat = proto.Seat
}
