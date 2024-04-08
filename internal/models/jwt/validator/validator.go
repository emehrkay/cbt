package validator

import (
	"context"
	"crypto"
	"fmt"
	"os"

	trainCTX "github.com/emehrkay/cbt/internal/models/jwt/ctx"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type JWTValidator interface {
	GetToken(tokenString string) (*jwt.Token, error)
	UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error)
}

func New(publicKeyFilePath string) (*validator, error) {
	keyBytes, err := os.ReadFile(publicKeyFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read public key file -- %w", err)
	}

	key, err := jwt.ParseEdPublicKeyFromPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf(`unalbe to get key from file: %v -- %w`, publicKeyFilePath, err)
	}
	val := &validator{
		key: key,
	}

	return val, nil
}

type validator struct {
	key crypto.PublicKey
}

func (v *validator) GetToken(tokenString string) (*jwt.Token, error) {
	var parse = func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("invalid token signer -- %v", token.Header["alg"])
		}

		return v.key, nil
	}

	token, err := jwt.Parse(tokenString, parse)
	if err != nil {
		return nil, fmt.Errorf("unable to parse token string -- %w", err)
	}

	return token, nil
}

func (v *validator) UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.New(codes.Unauthenticated, "no auth provided").Err()
	}
	tokens := headers.Get("jwt")
	if len(tokens) > 0 {
		tokenString := tokens[0]
		ctx = trainCTX.ContextWithToken(ctx, tokenString)
	}

	return handler(ctx, req)
}
