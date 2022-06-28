package db

import (
	"os"
	"testing"

	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/stkali/garden/token"
	"github.com/stkali/garden/util"
)
const (
	DSN = "postgres://root:password@localhost:5432/garden?sslmode=disable"
	DN  = "postgres"
)

var query Querier
var maker token.Maker
var duration = time.Hour * 24 * 14

func TestMain(m *testing.M) {
	conn, err := sql.Open(DN, DSN)
	util.CheckError("failed to create db connect", err)
	query = NewStore(conn)
	maker, err = token.NewMaker(token.GenerateSymmetricKey(), "jwt")
	util.CheckError("failed to create token maker", err)
	os.Exit(m.Run())
}
