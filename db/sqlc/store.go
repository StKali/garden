package db

import "database/sql"

func NewStore(db *sql.DB) Querier {
	return New(db)
}
