package storage

import (
	"database/sql"
	"fmt"
	"net/http"
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

type StoredRequest struct {
	RequestId string
	Response  *http.Response
	Finished  bool
}

type Storage interface {
	StoreRequest(blocking bool, method, url string) (*StoredRequest, error)
	StoreResponse(requestId string, response *http.Response) (*StoredRequest, error)
	GetRequest(requestId string) (*StoredRequest, error)
	Close() error
}

func GenerateStorage(name StorageName, postgresPassword string) (Storage, error) {
	if name == Postgres {
		var dbInfo string
		if postgresPassword != "" {
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
