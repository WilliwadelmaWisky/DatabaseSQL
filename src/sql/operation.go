package sql

type Operation interface {
	Call(database *Database) []byte
}
