package sql

import (
	"fmt"
	"slices"
)

type Database struct {
	rootPath string
	tables   []*Table
}

func NewDatabase(rootPath string, tables ...*Table) *Database {
	return &Database{
		rootPath: rootPath,
		tables:   tables,
	}
}

func (database *Database) Get(tableName string) (*Table, error) {
	index := slices.IndexFunc(database.tables, func(t *Table) bool { return t.Name == tableName })
	if index == -1 {
		return nil, fmt.Errorf("table not found: %s", tableName)
	}

	return database.tables[index], nil
}

func (database *Database) Create(tableName string, data []ColData) {
	columns := []*Column{}
	for _, colData := range data {
		columns = append(columns, &Column{
			Name:   colData.ColName,
			Type:   colData.ColType,
			Values: []string{},
		})
	}

	database.tables = append(database.tables, &Table{
		Name:    tableName,
		Columns: columns,
	})
}

func (database *Database) Save() {

}

func (database *Database) Load() {
	fmt.Printf("Path: %s\n", database.rootPath)
}
