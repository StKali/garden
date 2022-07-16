package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
	"github.com/stkali/garden/util"
)


func eqaulPayload(expect, actual *Payload) bool {
	
	if expect.ID != actual.ID {
		return false
	}
	if expect.Username != actual.Username{
		return false
	}
	if !expect.ExpiredAt.Equal(actual.ExpiredAt) {
		return false
	}
	return expect.IssuedAt.Equal(actual.IssuedAt)
}

func TestJWTMaker(t *testing.T) {

	maker, err := NewJWTMaker(GenerateSymmetricKey())
	require.NoError(t, err)

	username := util.RandInternalString(4, 16)
	duration := time.Minute

	token, createdPayload, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, createdPayload)

	parsedPayload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.True(t, eqaulPayload(createdPayload, parsedPayload))
}

func TestExpiredJWTToken(t *testing.T) {
	
	maker, err := NewJWTMaker(GenerateSymmetricKey())
	require.NoError(t, err)

	token, createdPayload, err := maker.CreateToken(util.RandInternalString(4, 16), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, createdPayload)
	
	parsedPayload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ExpiredToken.Error())
	require.Nil(t, parsedPayload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	createdPayload, err := NewPayload(util.RandInternalString(8, 16), time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, createdPayload)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, createdPayload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	maker, err := NewJWTMaker(util.RandInternalString(16, 32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	parsedPayload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, InvalidToken.Error())
	require.Nil(t, parsedPayload)
}
