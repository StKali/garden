package token

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	return &JWTMaker{secretKey: secretKey}, nil
}

func (j *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	//TODO implement me
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(j.secretKey))
}

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
