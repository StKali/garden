package db

import (
	"context"
	"log"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stkali/garden/util"
	"github.com/stretchr/testify/require"
)

func TestNewStore(t *testing.T) {
	require.NotNil(t, query)
}

func RandUser() User {

	password := util.RandInternalString(8, 36)
	hashPassword, err := util.HashPassword(password)
	if err != nil {
		log.Fatal("failed to create rand user")
	}
	return User{
		Username:          util.RandInternalString(4, 16),
		FullName:          util.RandInternalString(6, 18),
		HashedPassword:    hashPassword,
		Email:             util.RandEmail(),
		CreatedAt:         time.Now(),
		PasswordChangedAt: time.Now(),
	}
}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	_ = MakeAndSaveUser(ctx, t)
}

func CompareUser(expect, actual User) bool {
	if expect.Username != actual.Username {
		return false
	}
	if expect.FullName != actual.FullName {
		return false
	}
	if expect.Email != actual.Email {
		return false
	}
	if expect.HashedPassword != actual.HashedPassword {
		return false
	}
	return true
}

func MakeAndSaveUser(ctx context.Context, t *testing.T) User {
	user := RandUser()
	params := CreateUserParams{
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
		FullName:       user.FullName,
		Email:          user.Email,
	}
	createdUser, err := query.CreateUser(ctx, params)
	require.NoError(t, err)
	require.True(t, CompareUser(user, createdUser))
	return user
}


func TestGetUser(t *testing.T) {
	ctx := context.Background()
	user := MakeAndSaveUser(ctx, t)
	getUser, err := query.GetUser(ctx, user.Username)
	require.NoError(t, err)
	user.CreatedAt = getUser.CreatedAt
	user.PasswordChangedAt = getUser.PasswordChangedAt
	require.True(t, CompareUser(user, getUser))
}
