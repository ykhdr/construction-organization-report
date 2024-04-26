package db

import (
	"construction-organization-report/pkg/config"
	"database/sql"
)

func NewDB(dbConfig *config.DBConfig) (*sql.DB, error) {

	db, err := sql.Open("postgres", dbConfig.ConnectionInfo())
	if err != nil {
		return nil, err
	}

	return db, nil
}
