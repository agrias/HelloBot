package db

type Database interface {
	Get(string) ([]byte, error)
	Put(string, []byte) (error)
	Delete(string) (error)
	Keys() ([]string)
}