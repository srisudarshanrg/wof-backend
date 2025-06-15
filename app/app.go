package app

import "database/sql"

type Application struct {
	FrontendLink string
	DatabaseDSN  string
	Port         string
	DB           *sql.DB
	Deployed     bool
}
