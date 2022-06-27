package auth

import (
	"database/sql"
)

type Provider string

const (
	GoogleProvider Provider = "google"
)

type OAuthEntry struct {
	Provider Provider
	Subject  string
	UserId   string
}

type Repo interface {
	GetOAuthEntry(provider Provider, subject string) (*OAuthEntry, error)
	CreateOAuthEntry(provider Provider, subject, userId string) error
	CreateUser(user *User) (string, error)
	GetUser(id string) (*User, error)
}

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}

func (pr *PostgresRepo) GetOAuthEntry(provider Provider, subject string) (*OAuthEntry, error) {
	row := pr.db.QueryRow("SELECT provider, openid_id, user_id FROM \"openid\" WHERE provider=$1 AND openid_id=$2", provider, subject)
	if row == nil {
		return nil, sql.ErrNoRows
	}
	result := OAuthEntry{}
	if err := row.Scan(&result.Provider, &result.Subject, &result.UserId); err != nil {
		return nil, err
	}
	return &result, nil
}

func (pr *PostgresRepo) CreateOAuthEntry(provider Provider, subject, userId string) error {
	_, err := pr.db.Exec("INSERT INTO \"openid\"(provider, openid_id, user_id) VALUES($1,$2,$3)", provider, subject, userId)
	if err != nil {
		return err
	}
	return nil
}

func (pr *PostgresRepo) CreateUser(user *User) (string, error) {
	row := pr.db.QueryRow("INSERT INTO \"user\"(name) VALUES($1) RETURNING user_id", user.Name)
	if row == nil {
		return "", sql.ErrNoRows
	}
	var id string
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func (pr *PostgresRepo) GetUser(id string) (*User, error) {
	row := pr.db.QueryRow("SELECT user_id, name FROM \"user\" WHERE user_id=$1", id)
	if row == nil {
		return nil, sql.ErrNoRows
	}
	user := User{}
	if err := row.Scan(&user.Id, &user.Name); err != nil {
		return nil, err
	}
	return &user, nil
}
