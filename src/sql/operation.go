package sql

import (
	"encoding/json"
)

// Base contract of an sql operation
type Operation interface {
	Call(database *Database) ([]byte, error)
}

// Sql create operation, for creating new tables in the database
type CreateOperation struct {
	TableName string
	Data      []ColData
}

// Create operation execute method, creates a new table with table_name in the database with the data in the operation
func (operation *CreateOperation) Call(database *Database) ([]byte, error) {
	database.Create(operation.TableName, operation.Data)
	return nil, nil
}

// Sql insert operation, for inserting data to existing tables
type InsertOperation struct {
	TableName string
	Data      []RowData
}

// Insert operation execute method, inserts data in the operation to a table by the table_name in the operation
func (operation *InsertOperation) Call(database *Database) ([]byte, error) {
	table, err := database.Get(operation.TableName)
	if err != nil {
		return nil, err
	}

	err = table.Insert(operation.Data)
	return nil, err
}

// Sql select operation, for fetching data from the database
type SelectOperation struct {
	TableName   string
	ColumnNames []string
	Filters     []*Filter
	Sorters     []*Sorter
}

// Select operation execute method, fetches data from a table by table_name
func (operation *SelectOperation) Call(database *Database) ([]byte, error) {
	table, err := database.Get(operation.TableName)
	if err != nil {
		return nil, err
	}

	data, err := table.Get(operation.ColumnNames, operation.Filters)
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// Sql update operation, for updating values in existing tables
type UpdateOperation struct {
	TableName string
	Data      []RowData
	Filters   []*Filter
}

// Update operation execute method, updates row of a table by table_name
func (operation *UpdateOperation) Call(database *Database) ([]byte, error) {
	table, err := database.Get(operation.TableName)
	if err != nil {
		return nil, err
	}

	err = table.Update(operation.Data, operation.Filters)
	return nil, err
}

// Sql delete operation, for deleting rows in existing tables
type DeleteOperation struct {
	TableName string
	Filters   []*Filter
}

// Delete operation execute method, deletes rows included in the filters from a table by table_name
func (operation *DeleteOperation) Call(database *Database) ([]byte, error) {
	table, err := database.Get(operation.TableName)
	if err != nil {
		return nil, err
	}

	err = table.Delete(operation.Filters)
	return nil, err
}
