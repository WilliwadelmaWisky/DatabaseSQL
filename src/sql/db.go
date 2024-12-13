package sql

import (
	"fmt"
	"slices"
)

type Database struct {
	Tables []Table
}

func (database *Database) GetTable(tableName string) (*Table, error) {
	index := slices.IndexFunc(database.Tables, func(t Table) bool {
		return t.Name == tableName
	})

	if index == -1 {
		return nil, fmt.Errorf("table [%s] not found!", tableName)
	}

	return &database.Tables[index], nil
}

type Table struct {
	Name    string   `json:"table"`
	Columns []Column `json:"columns"`
}

func (table *Table) Get(match func(*Column) bool) []*Object {
	rowCount := len(table.Columns)
	colCount := len(table.Columns[0].Values)
	values := make([]*Object, rowCount)

	for row := 0; row < rowCount; row++ {
		obj := &Object{
			Values: make([]string, colCount),
		}

		for col := 0; col < colCount; col++ {
			obj.Values[col] = table.Columns[col].Values[row]
		}

		values[row] = obj
	}

	return values
}

func (table *Table) Insert() {

}

type Column struct {
	Name   string   `json:"column"`
	Type   string   `json:"type"`
	Values []string `json:"values"`
}

type Object struct {
	Values []string
}
