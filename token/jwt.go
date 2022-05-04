package token

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTMaker struct {
	secretKey string
}

// NewJWTMaker create a JWT token maker follow the Maker interface
func NewJWTMaker(secretKey string) (Maker, error) {
	return &JWTMaker{secretKey: secretKey}, nil
}

// CreateToken returns token, paylaod if creates token that specified username and duration success else error
func (j *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	//TODO implement me
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil, err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(j.secretKey))
	return token, payload, err
}

// VerifyToken return the payload if the token is valid passed else error 
func (j *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, InvalidToken
		}
		return []byte(j.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		err2, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(err2.Inner, ExpiredToken) {
			return nil, ExpiredToken
		}
		return nil, InvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, InvalidToken
	}

	return payload, nil
}

var _ Maker = (*JWTMaker)(nil)
