package issuer

import (
	"crypto"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTIssuer interface {
	Token(user string, roles []string) (string, error)
}

func New(privateKeyFilePath string, host string) (*JWT, error) {
	keyBytes, err := os.ReadFile(privateKeyFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read public key file -- %w", err)
	}

	key, err := jwt.ParseEdPrivateKeyFromPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf(`unalbe to get private key from file: %v -- %w`, privateKeyFilePath, err)
	}

	return &JWT{
		host: host,
		key:  key,
	}, nil
}

type JWT struct {
	host string
	key  crypto.PrivateKey
}

func (j *JWT) Token(userID string, roles []string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(&jwt.SigningMethodEd25519{}, jwt.MapClaims{
		"aud":   "api",
		"nbf":   now.Unix(),
		"iat":   now.Unix(),
		"exp":   now.Add(time.Minute).Unix(),
		"iss":   j.host,
		"user":  userID,
		"roles": roles,
	})

	tokenString, err := token.SignedString(j.key)
	if err != nil {
		return "", fmt.Errorf("unable to sign token -- %w", err)
	}

	return tokenString, nil
}
