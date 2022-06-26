package auth

import "database/sql"

type Repo interface {
}

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}
