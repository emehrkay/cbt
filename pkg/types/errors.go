package types

import "errors"

var (
	ErrInvalidAccess          error = errors.New("invalid access")
	ErrInvalidAccessAdminOnly error = errors.New("invalid access admin only")
	ErrNotFound               error = errors.New("not found")
	ErrSeatTaken              error = errors.New("seat taken")
	ErrTrainFull              error = errors.New("train full")
	ErrInvalidCar             error = errors.New("invalid car")
	ErrInvalidSeat            error = errors.New("invalid seat")
	ErrInvalidToken           error = errors.New("invalid token")
	ErrTokenNotFound          error = errors.New("token not found")
	ErrInvalidUser            error = errors.New("invalid user")
)
