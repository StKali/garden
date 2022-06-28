package db

import (
	"context"
	"log"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stkali/garden/util"
	"github.com/stretchr/testify/require"
)


func RandSession() Session {
	user := RandUser()
	token, payload, err := maker.CreateToken(user.Username, duration)
	if err != nil {
		log.Fatalf("failed to create user token, err: %s", err)
	}
	session := Session{
		ID           :uuid.New(),
		Username     :user.Username,
		RefreshToken :token,
		UserAgent    :util.RandInternalString(10, 20),
		ClientIp     :util.RandIP(),
		ExpiresAt    :payload.ExpiredAt,
		CreatedAt    :payload.IssuedAt,
	}
	return session
}

func MakeAndSaveSession(ctx context.Context, t *testing.T) Session {
	user := MakeAndSaveUser(ctx, t)
	token, payload, err := maker.CreateToken(user.Username, duration)
	require.NoError(t, err)
	session := Session{
		ID           :uuid.New(),
		Username     :user.Username,
		RefreshToken :token,
		UserAgent    :util.RandInternalString(10, 20),
		ClientIp     :util.RandIP(),
		ExpiresAt    :payload.ExpiredAt,
		CreatedAt    :payload.IssuedAt,
	}

	params := CreateSessionParams{
		ID           : session.ID,
		Username     : session.Username,
		RefreshToken : session.RefreshToken,
		UserAgent    : session.UserAgent,  
		ClientIp     : session.ClientIp,
		ExpiresAt    : session.ExpiresAt,
	}
	
	createdSession, err := query.CreateSession(ctx, params)
	require.NoError(t, err)
	session.CreatedAt = createdSession.CreatedAt
	session.ExpiresAt = createdSession.ExpiresAt
	require.Equal(t, session, createdSession)
	return session
}

func CompareSession(expect, actual Session) bool {
	if expect.ID != actual.ID {
		return false
	}
	if expect.Username != actual.Username {
		return false
	}
	if expect.RefreshToken != actual.RefreshToken {
		return false
	}
	if expect.UserAgent != actual.UserAgent {
		return false
	}
	if expect.IsBlocked != actual.IsBlocked {
		return false
	}
	return true
}

func TestCreateSession(t *testing.T) {
	ctx := context.Background()
	_ = MakeAndSaveSession(ctx, t)
}

func TestGetSession(t *testing.T) {
	ctx := context.Background()
	session := MakeAndSaveSession(ctx, t)
	getSession, err := query.GetSession(ctx, session.ID)
	require.NoError(t, err)
	require.True(t, CompareSession(session, getSession))
}
