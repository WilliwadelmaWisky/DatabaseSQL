package sql

import (
	"fmt"
	"slices"
)

// Represents a single databse
type Database struct {
	rootPath string
	tables   []*Table
}

// Create a new database.
//   - rootPath defines the location of the saved files
func NewDatabase(rootPath string, tables ...*Table) *Database {
	return &Database{
		rootPath: rootPath,
		tables:   tables,
	}
}

// Get a table from the database by name
func (database *Database) Get(tableName string) (*Table, error) {
	index := slices.IndexFunc(database.tables, func(t *Table) bool { return t.Name == tableName })
	if index == -1 {
		return nil, fmt.Errorf("table not found: %s", tableName)
	}

	return database.tables[index], nil
}

// Create a new empty table in the database
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

// Write database to disk
func (database *Database) Save() {

}

// Read database from disk
func (database *Database) Load() {
	fmt.Printf("Path: %s\n", database.rootPath)
}
