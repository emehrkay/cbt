package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/emehrkay/cbt/internal/models/cdc"
	trainCTX "github.com/emehrkay/cbt/internal/models/jwt/ctx"
	"github.com/emehrkay/cbt/internal/models/jwt/issuer"
	"github.com/emehrkay/cbt/internal/models/jwt/validator"
	"github.com/emehrkay/cbt/internal/storage"
	"github.com/emehrkay/cbt/internal/storage/demo"
	"github.com/emehrkay/cbt/pkg/types"
	"github.com/golang-jwt/jwt"
)

func NewFromEnvironment() (*Train, error) {
	var err error

	rpcPort := os.Getenv("RPC_PORT")
	if rpcPort == "" {
		rpcPort = ":7676"
	}

	store := demo.DemoStorage{}

	privateKeyFile := os.Getenv("PRIVATE_KEY_FILE")
	if privateKeyFile == "" {
		privateKeyFile = "./train/key"
	}

	privateKeyFile, err = filepath.Abs(privateKeyFile)
	if err != nil {
		panic(err)
	}

	publicKeyFile := os.Getenv("PUBLIC_KEY_FILE")
	if publicKeyFile == "" {
		publicKeyFile = "./train/key.pub"
	}

	publicKeyFile, err = filepath.Abs(publicKeyFile)
	if err != nil {
		panic(err)
	}

	jwtIssuer, err := issuer.New(privateKeyFile, rpcPort)
	if err != nil {
		panic(err)
	}

	jwtValidator, err := validator.New(publicKeyFile)
	if err != nil {
		panic(err)
	}

	capture := cdc.New()

	return New(store, jwtIssuer, jwtValidator, capture)
}

func New(store storage.Storage, jwtIssuer issuer.JWTIssuer, jwtValidator validator.JWTValidator, capture cdc.CDC) (*Train, error) {
	t := &Train{
		store:        store,
		jwtIssuer:    jwtIssuer,
		jwtValidator: jwtValidator,
		capture:      capture,
	}

	return t, nil
}

type Train struct {
	store        storage.Storage
	jwtIssuer    issuer.JWTIssuer
	jwtValidator validator.JWTValidator
	capture      cdc.CDC
}

func (t *Train) GetValidator() validator.JWTValidator {
	return t.jwtValidator
}

func (t *Train) newAuth(ctx context.Context) (*auth, error) {
	tokenString, err := trainCTX.ContextGetToken(ctx)
	if err != nil {
		return nil, err
	}

	token, err := t.jwtValidator.GetToken(*tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf(`unable to convert token claim to map -- %w`, types.ErrInvalidToken)
	}

	userID, ok := claims["user"].(string)
	if !ok {
		return nil, fmt.Errorf(`unable to extract user from claims -- %w`, types.ErrInvalidUser)
	}

	user, err := t.store.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return &auth{
		User:  *user,
		train: t,
	}, nil
}

type auth struct {
	train *Train
	User  types.User
}

func (a *auth) IsAdmin() bool {
	return a.User.Level == types.LevelAdmin
}

func (a *auth) CanEditTicket(ticket types.Ticket) (bool, error) {
	if a.IsAdmin() {
		return true, nil
	}

	user, err := a.train.store.GetUserByID(ticket.UserID)
	if err != nil {
		return false, fmt.Errorf(`unable to get ticket user -- %w`, err)
	}

	return user.ID == a.User.ID, nil
}
