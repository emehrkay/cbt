package service

import (
	"context"

	"github.com/emehrkay/cbt/pkg/types"
)

func (t *Train) Login(ctx context.Context, email, password string) (*types.LoginResponse, error) {
	user, err := t.store.Login(email, password)
	if err != nil {
		return nil, err
	}

	roles := []string{string(user.Level)}
	jwtToken, err := t.jwtIssuer.Token(user.ID, roles)
	if err != nil {
		return nil, err
	}

	resp := types.LoginResponse{
		JWT: jwtToken,
	}

	t.capture.Change(*user, "logged in")

	return &resp, nil
}
