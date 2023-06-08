package dbrepo

import (
	"database/sql"

	"github.com/TheDevCarnage/FortSmythesMotel/internals/config"
	"github.com/TheDevCarnage/FortSmythesMotel/internals/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB *sql.DB
}


func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo{
	return &postgresDBRepo{
		DB: conn,
		App: a,
	}
}