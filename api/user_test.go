package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	db "github.com/stkali/garden/db/sqlc"
	"github.com/stkali/garden/token"
	"github.com/stkali/garden/util"
	"github.com/stretchr/testify/require"
)

const (
	DSN = "postgres://root:password@localhost:5432/garden?sslmode=disable"
	DN  = "postgres"
)

var (
	query db.Querier
	maker token.Maker
)

func TestMain(m *testing.M) {
	conn, err := sql.Open(DN, DSN)
	util.CheckError("failed to create db connect", err)
	query = db.NewStore(conn)

	maker, err = token.NewMaker(token.GenerateSymmetricKey(), "poseto")
	util.CheckError("failed to create token maker", err)
	os.Exit(m.Run())
}

func generateRandUser(t *testing.T, password string) *db.User {
	hashPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	current := time.Now()
	return &db.User{
		Username:          util.RandInternalString(4, 32),
		FullName:          util.RandInternalString(4, 32),
		Email:             util.RandEmail(),
		HashedPassword:    hashPassword,
		PasswordChangedAt: current,
		CreatedAt:         current,
	}
}

func matchUser(t *testing.T, actUser *UserInfo, wantUser *db.User) {
	require.Equal(t, actUser.Username, wantUser.Username)
	require.Equal(t, actUser.FullName, wantUser.FullName)
	require.Equal(t, actUser.Email, wantUser.Email)
	require.NotEmpty(t, actUser.CreatedAt)
}

func TestCreateUser(t *testing.T) {
	password := util.RandString(8)
	user := generateRandUser(t, password)
	server := NewServer(query, maker)
	cases := []struct {
		Name  string
		Body  CreateUserRequest
		Check func(recorder *httptest.ResponseRecorder)
	}{
		{
			"OK",
			CreateUserRequest{
				Username:        user.Username,
				FullName:        user.FullName,
				Password:        password,
				ConfirmPassword: password,
				Email:           user.Email,
			},
			func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)
				res := new(UserInfo)
				err = json.Unmarshal(data, res)
				require.NoError(t, err)

				matchUser(t, res, user)
			},
		},
		{
			"DuplicateUsername",
			CreateUserRequest{
				Username:        user.Username,
				FullName:        user.FullName,
				Password:        password,
				ConfirmPassword: password,
				Email:           user.Email,
			},
			func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			"ConflictPassword",
			CreateUserRequest{
				Username:        util.RandInternalString(6, 16),
				FullName:        user.FullName,
				Password:        password,
				ConfirmPassword: password + "a",
				Email:           user.Email,
			},
			func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			"BadRequest",
			CreateUserRequest{
				Username:        util.RandInternalString(6, 16),
				Password:        password,
				ConfirmPassword: password,
			},
			func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			data, err := json.Marshal(c.Body)
			require.NoError(t, err)
			request := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader(data))
			server.engine.ServeHTTP(recorder, request)
			c.Check(recorder)
		})
	}
}

func TestLogin(t *testing.T) {

	password := util.RandInternalString(8, 32)
	hashPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	// create a valid user
	user, err := query.CreateUser(context.Background(), db.CreateUserParams{
		Username:       util.RandInternalString(8, 16),
		HashedPassword: hashPassword,
		FullName:       util.RandString(8),
		Email:          util.RandEmail(),
	})
	require.NoError(t, err)

	cases := []struct {
		Name  string
		Body  LoginRequest
		Check func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			"OK",
			LoginRequest{
				Username: user.Username,
				Password: password,
			},
			func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)
				res := new(UserInfo)
				err = json.Unmarshal(data, res)
				require.NoError(t, err)
				matchUser(t, res, &user)
			},
		},
		{
			"NotRegistered",
			LoginRequest{
				Username: user.Username + "a",
				Password: password,
			},
			func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			"InvalidPassword",
			LoginRequest{
				Username: user.Username,
				Password: password + "1",
			},
			func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			data, err := json.Marshal(c.Body)
			require.NoError(t, err)
			request := httptest.NewRequest(http.MethodGet, "/login", bytes.NewReader(data))
			recorder := httptest.NewRecorder()
			server := NewServer(query, maker)
			server.engine.ServeHTTP(recorder, request)
			c.Check(t, recorder)
		})
	}
}
