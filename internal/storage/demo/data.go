package demo

import "github.com/emehrkay/cbt/pkg/types"

const (
	ADMIN_EMAIL string = "admin@admin.com"
	ADMIN_PASS  string = "admin"
	USER_EMAIL  string = "user@user.com"
	USER_PASS   string = "user"
	USER2_EMAIL string = "user2@user.com"
	USER2_PASS  string = "user2"
	TRAIN_FROM  string = "London"
	TRAIN_TO    string = "France"
)

var (
	DemoTicketPrice types.Currency = types.Currency(20 * types.Dollar)
	maxSeatCount    int            = 10
	demoUsers       types.UserSet
	demoTrain       map[types.TrainCar]types.TicketSet
	train           types.Train
)

func init() {
	train = types.Train{
		Name:        "demo",
		From:        TRAIN_FROM,
		To:          TRAIN_TO,
		TicketPrice: DemoTicketPrice,
	}
	demoTrain = map[types.TrainCar]types.TicketSet{
		types.CarA: {},
		types.CarB: {},
	}

	demoUsers = types.UserSet{
		// user
		{
			ID:        "1",
			FirstName: "user first name",
			LastName:  "user last name",
			Email:     USER_EMAIL,
			Password:  USER_PASS,
			Level:     types.LevelUser,
		},
		{
			ID:        "2",
			FirstName: "user2_fname",
			LastName:  "user2_lname",
			Email:     USER2_EMAIL,
			Password:  USER2_PASS,
			Level:     types.LevelUser,
		},

		// admin
		{
			ID:        "3",
			FirstName: "admin first name",
			LastName:  "admin last name",
			Email:     ADMIN_EMAIL,
			Password:  ADMIN_PASS,
			Level:     types.LevelAdmin,
		},
	}
}
