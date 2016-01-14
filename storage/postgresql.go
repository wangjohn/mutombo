package storage

import (
	"database/sql"
	"net/http"
)

type PostgresStorage struct {
	DB *sql.DB
}

func (s PostgresStorage) StoreRequest(blocking bool, method, url string) (*StoredRequest, error) {

}

func (s PostgresStorage) StoreResponse(requestId string, response *http.Response) (*StoredRequest, error) {

}

func (s PostgresStorage) GetRequest(requestId string) (*StoredRequest, error) {

}

func (s PostgresStorage) Close() error {
	return s.DB.Close()
}
