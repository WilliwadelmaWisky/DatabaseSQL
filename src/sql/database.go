package sql

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
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
func (database *Database) Create(tableName string, data []ColData) error {
	index := slices.IndexFunc(database.tables, func(t *Table) bool { return t.Name == tableName })
	if index != -1 {
		return fmt.Errorf("table already exists: %s", tableName)
	}

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

	return nil
}

// Delete a table from the database
func (database *Database) Delete(tableName string) error {
	index := slices.IndexFunc(database.tables, func(t *Table) bool { return t.Name == tableName })
	if index == -1 {
		return fmt.Errorf("table not found: %s", tableName)
	}

	database.tables = slices.Delete(database.tables, index, index+1)
	return nil
}

// Write database to disk
func (database *Database) Save() error {
	fmt.Printf("Save: %s\n", database.rootPath)

	files, err := os.ReadDir(database.rootPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := filepath.Join(database.rootPath, file.Name())
		if !strings.HasSuffix(filePath, ".json") {
			continue
		}

		os.Remove(filePath)
	}

	for _, table := range database.tables {
		filePath := filepath.Join(database.rootPath, fmt.Sprintf("%s.json", table.Name))
		data, err := json.Marshal(table)
		if err != nil {
			return err
		}

		err = os.WriteFile(filePath, data, fs.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

// Read database from disk
func (database *Database) Load() error {
	fmt.Printf("Load: %s\n", database.rootPath)

	files, err := os.ReadDir(database.rootPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := filepath.Join(database.rootPath, file.Name())
		if !strings.HasSuffix(filePath, ".json") {
			continue
		}

		data, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		table := &Table{}
		err = json.Unmarshal(data, table)
		if err != nil {
			return err
		}

		database.tables = append(database.tables, table)
	}

	return nil
}
