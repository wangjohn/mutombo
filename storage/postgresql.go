package storage

import (
	"bytes"
	"database/sql"
	"io/ioutil"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	storeRequestQuery  = `INSERT INTO requests(blocking, method, url) VALUES ($1, $2, $3) returning id;`
	storeResponseQuery = `UPDATE requests SET body = $1 WHERE id = $2;`
	//storeHeadersQuery = `INSERT INTO headers(request_id, name, value) VALUES ($1, $2, $3);`
	getRequestQuery = `SELECT finished, body, status_code FROM requests WHERE id = $1;`
	getHeadersQuery = `SELECT name, value FROM headers WHERE request_id = $1;`
)

type PostgresStorage struct {
	DB *sql.DB
}

// Stores request into database for future lookup
func (s PostgresStorage) StoreRequest(blocking bool, method, url string) (*StoredRequest, error) {
	var requestId string
	err := s.DB.QueryRow(storeRequestQuery, blocking, method, url).Scan(&requestId)
	if err != nil {
		return nil, err
	}

	encryptedId, err := encryptRequestId(requestId)
	if err != nil {
		return nil, err
	}
	storedReq := StoredRequest{
		RequestId: encryptedId,
		Finished:  false,
	}
	return &storedReq, nil
}

// Stores the response of a request into database for lookup with a stored request
func (s PostgresStorage) StoreResponse(requestId string, response *http.Response) (*StoredRequest, error) {
	reqId, err := decryptRequestId(requestId)
	if err != nil {
		return nil, err
	}
	respByteBody := response.Body
	err = s.DB.QueryRow(storeResponseQuery, respByteBody, reqId).Scan()
	if err != nil {
		return nil, err
	}
	// TODO(wangjohn): Store headers
	return s.getStoredRequest(reqId)
}

// Returns the current state of a stored request from the requestId
func (s PostgresStorage) GetRequest(requestId string) (*StoredRequest, error) {
	reqId, err := decryptRequestId(requestId)
	if err != nil {
		return nil, err
	}
	return s.getStoredRequest(reqId)
}

func (s PostgresStorage) Close() error {
	return s.DB.Close()
}

// Helper function for returning a StoredRequest object for a request
func (s PostgresStorage) getStoredRequest(reqId string) (*StoredRequest, error) {
	var finished bool
	var body []byte
	var statusCode int
	err := s.DB.QueryRow(getRequestQuery, reqId).Scan(&finished, &body, &statusCode)
	if err != nil {
		return nil, err
	}

	respHeader, err := s.getHeader(reqId)
	if err != nil {
		return nil, err
	}

	bodyReader := ioutil.NopCloser(bytes.NewReader(body))
	httpResp := http.Response{
		StatusCode: statusCode,
		Header:     *respHeader,
		Body:       bodyReader,
	}

	storedReq := StoredRequest{
		RequestId: reqId,
		Response:  &httpResp,
		Finished:  finished,
	}
	return &storedReq, err
}

// Helper function for getting the stored header object for a request
func (s PostgresStorage) getHeader(reqId string) (*http.Header, error) {
	rows, err := s.DB.Query(getHeadersQuery, reqId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	header := http.Header{}
	for rows.Next() {
		var name, value string
		if err := rows.Scan(&name, &value); err != nil {
			return nil, err
		}
		header.Set(name, value)
	}
	return &header, nil
}
