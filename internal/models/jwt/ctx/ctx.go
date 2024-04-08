package ctx

import (
	"context"

	"github.com/emehrkay/cbt/pkg/types"
)

type JWTContextKeyType string

const JWTContextKey JWTContextKeyType = "jwt-ctx"

func ContextWithToken(ctx context.Context, tokenString string) context.Context {
	return context.WithValue(ctx, JWTContextKey, tokenString)
}

func ContextGetToken(ctx context.Context) (*string, error) {
	val := ctx.Value(JWTContextKey)
	if val == nil {
		return nil, types.ErrTokenNotFound
	}

	tokenString, ok := val.(string)
	if !ok {
		return nil, types.ErrInvalidToken
	}

	return &tokenString, nil
}
