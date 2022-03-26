package db


import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/stkali/garden/util"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

const (
	DSN = "postgres://root:password@localhost:5432/garden?sslmode=disable"
	DN  = "postgres"
)

var query Querier

func TestMain(m *testing.M) {
	conn, err := sql.Open(DN, DSN)
	if err != nil {
		log.Fatalf("failed to create db connect, err: %s", err)
	}
	query = NewStore(conn)
	os.Exit(m.Run())
}

func TestNewStore(t *testing.T) {
	require.NotNil(t, query)
}

func TestCreateUser(t *testing.T) {
	params := CreateUserParams{
		Username:       util.RandInternalString(0, 10),
		HashedPassword: util.RandInternalString(0, 10),
		FullName:       util.RandInternalString(0, 10),
		Email:          util.RandEmail(),
	}
	ctx := context.Background()
	user, err := query.CreateUser(ctx, params)
	require.NoError(t, err)
	require.Equal(t, user.Email, params.Email)
}

func TestGetUser(t *testing.T) {
	params := CreateUserParams{
		Username:       util.RandInternalString(0, 10),
		HashedPassword: util.RandInternalString(0, 10),
		FullName:       util.RandInternalString(0, 10),
		Email:          util.RandEmail(),
	}
	ctx := context.Background()
	user, err := query.CreateUser(ctx, params)
	require.NoError(t, err)
	require.Equal(t, user.Email, params.Email)

	user2, err := query.GetUser(ctx, params.Username)
	require.NoError(t, err)
	require.Equal(t, user2, user)
}

