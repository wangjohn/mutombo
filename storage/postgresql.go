package storage

import (
	"database/sql"
)

type PostgresStorage struct {
	DB *sql.DB
}

func (s PostgresStorage) Close() error {
	return s.DB.Close()
}
