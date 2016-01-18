package storage

import (
	"database/sql"
	"io/ioutil"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

const (
	storeRequestQuery  = `INSERT INTO requests(blocking, method, url) VALUES ($1, $2, $3) returning id;`
	storeResponseQuery = `UPDATE requests SET body = $1, finished = 't', status_code = $2 WHERE id = $3;`
	storeHeadersQuery  = `INSERT INTO headers(request_id, name, value) VALUES ($1, $2, $3);`
	getRequestQuery    = `SELECT finished, body, status_code FROM requests WHERE id = $1;`
	getHeadersQuery    = `SELECT name, value FROM headers WHERE request_id = $1;`
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
	respByteBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	// Store the response body and status code to the request id
	_, err = s.DB.Exec(storeResponseQuery, string(respByteBody), response.StatusCode, reqId)
	if err != nil {
		return nil, err
	}
	// Store headers to the request id
	stmt, err := s.DB.Prepare(storeHeadersQuery)
	if err != nil {
		return nil, err
	}
	for k, _ := range response.Header {
		_, err = stmt.Exec(reqId, k, response.Header.Get(k))
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
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
	var body string
	var sc sql.NullInt64
	err := s.DB.QueryRow(getRequestQuery, reqId).Scan(&finished, &body, &sc)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	statusCode := 0
	if sc.Valid {
		statusCode = int(sc.Int64)
	}

	respHeader, err := s.getHeader(reqId)
	if err != nil {
		return nil, err
	}

	bodyReader := ioutil.NopCloser(strings.NewReader(body))
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
