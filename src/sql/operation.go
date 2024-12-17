package sql

import (
	"encoding/json"
)

type Operation interface {
	Call(database *Database) ([]byte, error)
}

type CreateOperation struct {
	TableName string
	Data      []ColData
}

func (operation *CreateOperation) Call(database *Database) ([]byte, error) {
	database.Create(operation.TableName, operation.Data)
	return nil, nil
}

type InsertOperation struct {
	TableName string
	Data      []RowData
}

func (operation *InsertOperation) Call(database *Database) ([]byte, error) {
	table, err := database.Get(operation.TableName)
	if err != nil {
		return nil, err
	}

	table.Insert(operation.Data)
	return nil, nil
}

type SelectOperation struct {
	TableName   string
	ColumnNames []string
	Filters     []*Filter
	Sorters     []*Sorter
}

func (operation *SelectOperation) Call(database *Database) ([]byte, error) {
	table, err := database.Get(operation.TableName)
	if err != nil {
		return nil, err
	}

	data := table.Get(operation.ColumnNames, operation.Filters)
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

type UpdateOperation struct {
	TableName string
	Data      []RowData
	Filters   []*Filter
}

func (operation *UpdateOperation) Call(database *Database) ([]byte, error) {
	table, err := database.Get(operation.TableName)
	if err != nil {
		return nil, err
	}

	table.Update(operation.Data, operation.Filters)
	return nil, nil
}

type DeleteOperation struct {
	TableName string
	Filters   []*Filter
}

func (operation *DeleteOperation) Call(database *Database) ([]byte, error) {
	table, err := database.Get(operation.TableName)
	if err != nil {
		return nil, err
	}

	table.Delete(operation.Filters)
	return nil, nil
}
