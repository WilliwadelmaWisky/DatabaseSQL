package sql

import (
	"encoding/json"
)

type SelectOperation struct {
	ColumnNames []string
	TableName   string
	Filters     []*Filter
}

func (operation *SelectOperation) Call(database *Database) []byte {
	table, _ := database.Get(operation.TableName)
	objects := table.Get(operation.ColumnNames, operation.Filters)
	bytes, _ := json.Marshal(objects)
	return bytes
}
