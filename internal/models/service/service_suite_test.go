package service_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/emehrkay/cbt/internal/models/service"
	"github.com/emehrkay/cbt/internal/storage/demo"
	"github.com/emehrkay/cbt/pkg/types"
)

var (
	err         error
	testService *service.Train
)

var _ = BeforeSuite(func() {
	testService, err = service.NewFromEnvironment()
	Expect(err).To(BeNil())
})

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Layer Test Suite")
}

func UserLogin() *types.LoginResponse {
	loginResp, err := testService.Login(context.Background(), demo.USER_EMAIL, demo.USER_PASS)
	Expect(err).To(BeNil())

	return loginResp
}

func User2Login() *types.LoginResponse {
	loginResp, err := testService.Login(context.Background(), demo.USER2_EMAIL, demo.USER2_PASS)
	Expect(err).To(BeNil())

	return loginResp
}

func AdminLogin() *types.LoginResponse {
	loginResp, err := testService.Login(context.Background(), demo.ADMIN_EMAIL, demo.ADMIN_PASS)
	Expect(err).To(BeNil())

	return loginResp
}

func BuyTicket(ctx context.Context) *types.Ticket {
	ticketResp, err := testService.BuyTicket(ctx, types.BuyTicket{})
	Expect(err).To(BeNil())

	return ticketResp
}
