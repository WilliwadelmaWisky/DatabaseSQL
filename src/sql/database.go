package sql

import (
	"fmt"
	"slices"
)

// Represents a single databse
type Database struct {
	rootPath string   // Database location on the disk
	tables   []*Table // Database tables
}

// Create a new database.
//   - rootPath defines the location of the saved files
func NewDatabase(rootPath string, tables ...*Table) *Database {
	return &Database{
		rootPath: rootPath,
		tables:   tables,
	}
}

// Represents a metadata of the database
type InformationSchema struct {
	Tables []string `json:"tables"`
}

// Create a new information_schema
func NewInformationSchema(database *Database) *InformationSchema {
	return &InformationSchema{
		Tables: Map(database.tables, func(table *Table) string { return table.Name }),
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
	fmt.Printf("Save databse to %s\n", database.rootPath)
}

// Read database from disk
func (database *Database) Load() {
	fmt.Printf("Load database from %s\n", database.rootPath)
}
