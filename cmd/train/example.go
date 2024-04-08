package train

import (
	"context"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	trainCTX "github.com/emehrkay/cbt/internal/models/jwt/ctx"
	"github.com/emehrkay/cbt/internal/storage/demo"
	rpcBase "github.com/emehrkay/cbt/pkg/rpc"
	"github.com/emehrkay/cbt/pkg/types"
)

func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use:   "ex1",
		Short: "runs example 1",
		Long:  "logs in with a user, buys a ticket",
		Run: func(cmd *cobra.Command, args []string) {
			exampleOne()
		},
	})

	RootCmd.AddCommand(&cobra.Command{
		Use:   "ex2",
		Short: "runs example 2",
		Long:  "logs in with a user, buys a ticket. Then logs in with admin and looks at the train details",
		Run: func(cmd *cobra.Command, args []string) {
			exampleTwo()
		},
	})

	RootCmd.AddCommand(&cobra.Command{
		Use:   "ex3",
		Short: "runs example 3",
		Long:  "logs in with a user and attemps to view train details, an admin-only action",
		Run: func(cmd *cobra.Command, args []string) {
			exampleThree()
		},
	})

	RootCmd.AddCommand(&cobra.Command{
		Use:   "ex4",
		Short: "runs example 4",
		Long:  "logs in with a user and adds a ticket, logs in as admin and removes ticket from train",
		Run: func(cmd *cobra.Command, args []string) {
			exampleFour()
		},
	})

	RootCmd.AddCommand(&cobra.Command{
		Use:   "ex5",
		Short: "runs example 4",
		Long:  "logs in with a user and adds a ticket, logs in as admin and updates the ticket. the user logs in and changes the ticket again. finally another user tries to unsuccessfully edit the first user's ticket",
		Run: func(cmd *cobra.Command, args []string) {
			exampleFive()
		},
	})
}

func newClient(withInterceptor bool) (*grpc.ClientConn, rpcBase.TrainClient) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if withInterceptor {
		opts = append(opts, grpc.WithUnaryInterceptor(trainCTX.UnaryClientInterceptor))
	}

	conn, err := grpc.NewClient(rpcPort, opts...)
	if err != nil {
		panic(err)
	}

	return conn, rpcBase.NewTrainClient(conn)
}

func exampleOne() {
	// create connection
	startRPCServer()

	// create the client
	conn, client := newClient(false)
	ctx := context.Background()
	defer conn.Close()

	// login
	resp, err := client.Login(ctx, &rpcBase.LoginRequest{
		Email:    demo.USER_EMAIL,
		Password: demo.USER_PASS,
	})
	if err != nil {
		panic(err)
	}

	log(`logged in as "user" with jwt: %v`, resp.Jwt)

	conn, client = newClient(true)
	defer conn.Close()

	ctx = trainCTX.ContextWithToken(ctx, resp.Jwt)

	// get next available seat
	seats, err := client.TrainOpenSeats(ctx, nil)
	if err != nil {
		panic(err)
	}
	log("open seats \n\t%+v", seats)

	// purchase a ticket
	ticket, err := client.BuyTicket(ctx, &rpcBase.BuyTicketRequest{})
	if err != nil {
		panic(err)
	}
	log("bought ticket\n\t%+v", ticket)

	// check that the seat has been taken
	seats, err = client.TrainOpenSeats(ctx, nil)
	if err != nil {
		panic(err)
	}
	log("open seats left \n\t%+v", seats)
}

func exampleTwo() {
	// buy five tickets
	exampleOne()
	exampleOne()
	exampleOne()
	exampleOne()

	// create the client
	conn, client := newClient(false)
	ctx := context.Background()
	defer conn.Close()

	// login
	resp, err := client.Login(ctx, &rpcBase.LoginRequest{
		Email:    demo.ADMIN_EMAIL,
		Password: demo.ADMIN_PASS,
	})
	if err != nil {
		panic(err)
	}

	log(`logged in as "admin" with jwt: %v`, resp.Jwt)

	conn, client = newClient(true)
	defer conn.Close()

	ctx = trainCTX.ContextWithToken(ctx, resp.Jwt)

	// get all of the purchased tickets
	details, err := client.TrainDetails(ctx, nil)
	if err != nil {
		panic(err)
	}

	log("train details:\n\t%+v", details)
}

func exampleThree() {
	// create connection
	startRPCServer()

	// create the client
	conn, client := newClient(false)
	ctx := context.Background()
	defer conn.Close()

	// login
	resp, err := client.Login(ctx, &rpcBase.LoginRequest{
		Email:    demo.USER_EMAIL,
		Password: demo.USER_PASS,
	})
	if err != nil {
		panic(err)
	}

	log(`logged in as "user" with jwt: %v`, resp.Jwt)

	conn, client = newClient(true)
	defer conn.Close()
	ctx = trainCTX.ContextWithToken(ctx, resp.Jwt)

	//attempt to do admin-only action
	_, err = client.TrainDetails(ctx, nil)
	if err != nil {
		log("user attempted to do admin-only thing:\n\t%v", err.Error())
	}
}

func exampleFour() {
	exampleOne()

	// login as admin, delete ticket
	conn, client := newClient(false)
	ctx := context.Background()
	defer conn.Close()

	// login
	resp, err := client.Login(ctx, &rpcBase.LoginRequest{
		Email:    demo.ADMIN_EMAIL,
		Password: demo.ADMIN_PASS,
	})
	if err != nil {
		panic(err)
	}

	log(`logged in as "admin" with jwt: %v`, resp.Jwt)

	conn, client = newClient(true)
	defer conn.Close()

	ctx = trainCTX.ContextWithToken(ctx, resp.Jwt)

	// get all of the purchased tickets
	details, err := client.TrainDetails(ctx, nil)
	if err != nil {
		panic(err)
	}

	var ticketID string

	for _, ticket := range details.Trains {
		if ticket != nil {
			ticketID = ticket.Ticket.Id
			log("removing ticket:\n\t%+v", ticket)
		}
	}

	_, err = client.RemoveTicket(ctx, &rpcBase.RemoveTicketRequest{
		TicketId: ticketID,
	})
	if err != nil {
		panic(err)
	}

	// get all of the purchased tickets
	details, err = client.TrainDetails(ctx, nil)
	if err != nil {
		panic(err)
	}

	log("train details:\n\t%+v", details)
}

func exampleFive() {
	exampleOne()

	// login as admin, delete ticket
	conn, client := newClient(false)
	ctx := context.Background()
	defer conn.Close()

	// login
	resp, err := client.Login(ctx, &rpcBase.LoginRequest{
		Email:    demo.ADMIN_EMAIL,
		Password: demo.ADMIN_PASS,
	})
	if err != nil {
		panic(err)
	}

	log(`logged in as "admin" with jwt: %v`, resp.Jwt)

	conn, client = newClient(true)
	defer conn.Close()

	ctx = trainCTX.ContextWithToken(ctx, resp.Jwt)

	// get all of the purchased tickets
	details, err := client.TrainDetails(ctx, nil)
	if err != nil {
		panic(err)
	}

	var ticketID string
	car := types.TrainCar("b")
	seat := int32(4)

	for _, ticket := range details.Trains {
		if ticket != nil {
			ticketID = ticket.Ticket.Id
			log("changing ticket:\n\t%+v", ticket)
		}
	}

	updated, err := client.ChangeTicket(ctx, &rpcBase.ChangeTicketRequest{
		TicketId: ticketID,
		Car:      string(car),
		Seat:     seat,
	})
	if err != nil {
		panic(err)
	}

	log("admin changed ticket to:\n\t%+v", updated)

	// user
	resp, err = client.Login(ctx, &rpcBase.LoginRequest{
		Email:    demo.USER_EMAIL,
		Password: demo.USER_PASS,
	})
	if err != nil {
		panic(err)
	}

	ctx = trainCTX.ContextWithToken(ctx, resp.Jwt)

	car = types.CarA
	seat = int32(6)
	updated, err = client.ChangeTicket(ctx, &rpcBase.ChangeTicketRequest{
		TicketId: updated.Id,
		Car:      string(car),
		Seat:     seat,
	})
	if err != nil {
		panic(err)
	}

	log("user changed ticket to:\n\t%+v", updated)

	resp, err = client.Login(ctx, &rpcBase.LoginRequest{
		Email:    demo.USER2_EMAIL,
		Password: demo.USER2_PASS,
	})
	if err != nil {
		panic(err)
	}

	ctx = trainCTX.ContextWithToken(ctx, resp.Jwt)

	car = types.CarB
	seat = int32(67)
	_, err = client.ChangeTicket(ctx, &rpcBase.ChangeTicketRequest{
		TicketId: updated.Id,
		Car:      string(car),
		Seat:     seat,
	})

	log(`user2 cannot update another user's ticket\n\t%v`, err)
}
