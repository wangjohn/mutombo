package storage

type StorageName int

const (
	Postgres StorageName = iota
)

type Storage interface {
	Close() error
}
