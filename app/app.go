package app

import "database/sql"

type Application struct {
	ProductionFrontendLink     string
	DevelopmentFrontendLink    string
	ProductionFrontendLinkTest string
	DatabaseDSN                string
	Port                       string
	DB                         *sql.DB
	Deployed                   bool
}
