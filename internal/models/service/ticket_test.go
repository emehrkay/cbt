package service_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	trainCTX "github.com/emehrkay/cbt/internal/models/jwt/ctx"
	"github.com/emehrkay/cbt/pkg/types"
)

var _ = Describe("tickets", func() {
	Context("ChangeTicket", func() {
		Context("failures", func() {
			It("should require a jwt in the context", func() {
				_, err := testService.ChangeTicket(context.Background(), types.ChangeTicket{})
				Expect(err.Error()).To(Equal(types.ErrTokenNotFound.Error()))
			})

			It("should error if the ticketID is invalid", func() {
				ctx := context.Background()
				user := UserLogin()

				ctx = trainCTX.ContextWithToken(ctx, user.JWT)
				ticket := BuyTicket(ctx)

				user2 := User2Login()
				ctx = trainCTX.ContextWithToken(ctx, user2.JWT)

				_, err = testService.ChangeTicket(ctx, types.ChangeTicket{
					TicketID: ticket.ID + "xxxyyy",
					Seat:     9,
				})
				Expect(err.Error()).To(ContainSubstring(types.ErrNotFound.Error()))
			})

			It("should error if the user doesnt have access to change the ticket", func() {
				ctx := context.Background()
				user := UserLogin()

				ctx = trainCTX.ContextWithToken(ctx, user.JWT)
				ticket := BuyTicket(ctx)

				user2 := User2Login()
				ctx = trainCTX.ContextWithToken(ctx, user2.JWT)

				_, err = testService.ChangeTicket(ctx, types.ChangeTicket{
					TicketID: ticket.ID,
					Seat:     9,
				})
				Expect(err.Error()).To(ContainSubstring(types.ErrInvalidAccess.Error()))
			})
		})

		Context("successes", func() {
			It("should allow a user to update their own ticket", func() {
				ctx := context.Background()
				user := UserLogin()

				ctx = trainCTX.ContextWithToken(ctx, user.JWT)
				ticket := BuyTicket(ctx)

				newSeat := int32(9)
				updated, err := testService.ChangeTicket(ctx, types.ChangeTicket{
					TicketID: ticket.ID,
					Seat:     newSeat,
				})
				Expect(err).To(BeNil())
				Expect(updated.Seat).To(Equal(newSeat))
			})

			It("should allow an admin to update another user's ticket", func() {
				ctx := context.Background()
				user := UserLogin()

				ctx = trainCTX.ContextWithToken(ctx, user.JWT)
				ticket := BuyTicket(ctx)

				admin := AdminLogin()
				ctx = trainCTX.ContextWithToken(ctx, admin.JWT)

				newSeat := int32(9)
				updated, err := testService.ChangeTicket(ctx, types.ChangeTicket{
					TicketID: ticket.ID,
					Seat:     newSeat,
				})
				Expect(err).To(BeNil())
				Expect(updated.Seat).To(Equal(newSeat))
			})
		})
	})
})
