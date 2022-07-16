package token

import (
	"testing"
	"time"

	"github.com/stkali/garden/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(GenerateSymmetricKey())
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := util.RandInternalString(4, 8)
	duration := time.Minute

	token, createdPayload, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, createdPayload)
	parsedPayload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, parsedPayload)

	require.True(t, eqaulPayload(createdPayload, parsedPayload))
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(GenerateSymmetricKey())
	require.NoError(t, err)

	token, createdPayload, err := maker.CreateToken(util.RandInternalString(6, 12), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, createdPayload)

	parsedpayload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ExpiredToken.Error())
	require.Nil(t, parsedpayload)
}
