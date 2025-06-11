package app

import (
	"database/sql"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func (app *Application) ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("pgx", app.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	return db, nil
}
