package sql

import "encoding/json"

type Operation interface {
	Call(database *Database) []byte
}

type CreateOperation struct {
	TableName string
	Data      []ColData
}

func (operation *CreateOperation) Call(database *Database) []byte {
	database.Create(operation.TableName, operation.Data)
	return nil
}

type InsertOperation struct {
	TableName string
	Data      []ValData
}

func (operation *InsertOperation) Call(database *Database) []byte {
	table, _ := database.Get(operation.TableName)
	table.Insert(operation.Data)
	return nil
}

type SelectOperation struct {
	TableName   string
	ColumnNames []string
	Filters     []*Filter
	Sorters     []*Sorter
}

func (operation *SelectOperation) Call(database *Database) []byte {
	table, _ := database.Get(operation.TableName)
	objects := table.Get(operation.ColumnNames, operation.Filters)
	bytes, _ := json.Marshal(objects)
	return bytes
}

type UpdateOperation struct {
}

func (operation *UpdateOperation) Call(database *Database) []byte {
	return nil
}

type DeleteOperation struct {
	TableName string
	Filters   []*Filter
}

func (operation *DeleteOperation) Call(database *Database) []byte {
	table, _ := database.Get(operation.TableName)
	table.Delete(operation.Filters)
	return nil
}
