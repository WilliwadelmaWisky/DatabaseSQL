package sql

import (
	"encoding/json"
	"slices"
)

type SelectOperation struct {
	ColumnNames []string
	TableName   string
}

func (operation *SelectOperation) Call(database *Database) []byte {
	table, _ := database.Get(operation.TableName)
	objects := table.Get(func(column *Column) bool {
		if len(operation.ColumnNames) == 1 && operation.ColumnNames[0] == "*" {
			return true // SELECT ALL
		}

		return slices.ContainsFunc(operation.ColumnNames, func(columnName string) bool { return columnName == column.Name })
	})

	bytes, _ := json.Marshal(objects)
	return bytes
}
