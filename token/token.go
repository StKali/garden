package token

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stkali/garden/util"
	"golang.org/x/crypto/chacha20poly1305"
	"time"
)

var (
	InvalidToken = errors.New("invalid token")
	ExpiredToken = errors.New("expired token")
)

type Maker interface {
	CreateToken(username string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ExpiredToken
	}
	return nil
}

func GenerateSymmetricKey() string {
	return util.RandString(chacha20poly1305.KeySize)
}
