package dbrepo

import (
	"database/sql"

	"github.com/kiniconnect/booking-app/internals/config"
)

type PostgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresDBRepo(app *config.AppConfig, conn *sql.DB) *PostgresDBRepo {
	return &PostgresDBRepo{
		App: app,
		DB:  conn,
	}
}


