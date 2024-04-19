package security // import "github.com/amieldelatorre/notifi/service/security"

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func CreateAccessToken(claims UserClaims, signingKey []byte) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := accessToken.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ParseAccessToken(tokenString string, signingKey []byte) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if !token.Valid {
		return nil, &InvalidAccessToken{}
	}

	if err != nil {
		return nil, err
	}

	if token == nil || token.Claims == nil {
		return nil, errors.New("parsed access token or it's claims is nil")
	}

	value, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New("could not parse JWT Token as UserClaims")
	}

	return value, nil

}

type InvalidAccessToken struct{}

func (e *InvalidAccessToken) Error() string {
	return "Invalid access token"
}
