package storage

import (
	"database/sql"
	"fmt"
	"os"
)

type StorageName int

const (
	Postgres StorageName = iota
)

const (
	postgresDriverName = "postgres"
	postgresHost       = "localhost"
	postgresUser       = "postgres"
	postgresDBName     = "mutombo"
	postgresSSLMode    = "disable"
)

type Storage interface {
	Close() error
}

func GenerateStorage(name StorageName) (Storage, error) {
	if name == Postgres {
		var dbInfo string
		if os.Getenv("ENVIRONMENT") == "PROD" {
			postgresPassword := os.Getenv("POSTGRES_PASSWORD")
			dbInfo = fmt.Sprintf("host=%v user=%s password=%s dbname=%s sslmode=%s",
				postgresHost, postgresUser, postgresPassword, postgresDBName, postgresSSLMode)
		} else {
			dbInfo = fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s",
				postgresHost, postgresUser, postgresDBName, postgresSSLMode)
		}
		db, err := sql.Open(postgresDriverName, dbInfo)
		return PostgresStorage{db}, err
	} else {
		return nil, fmt.Errorf("Invalid storage mechanism: %s", name)
	}
}
