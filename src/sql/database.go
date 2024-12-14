package sql

import (
	"fmt"
	"slices"
)

type Database struct {
	Tables []*Table `json:"tables"`
}

func (database *Database) Get(tableName string) (*Table, error) {
	index := slices.IndexFunc(database.Tables, func(t *Table) bool {
		return t.Name == tableName
	})

	if index == -1 {
		return nil, fmt.Errorf("table [%s] not found", tableName)
	}

	return database.Tables[index], nil
}

func (database *Database) Create(tableName string, data []ColData) {
	columns := []*Column{}
	for _, colData := range data {
		columns = append(columns, &Column{
			Name: colData.Name,
			Type: colData.Type,
		})
	}

	table := &Table{Name: tableName, Columns: columns}
	database.Tables = append(database.Tables, table)
}

type Object struct {
	Values []string
}
