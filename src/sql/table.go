package sql

import (
	"fmt"
	"slices"
)

// Represents a single table in the database
type Table struct {
	Name    string    `json:"table"`   // Table name
	Columns []*Column `json:"columns"` // Table columns
}

// An object to return by get method
type TableData struct {
	Columns []string   `json:"columns"` // Table column name array
	Data    [][]string `json:"data"`    // Table data, array of rows
}

type RowData struct {
	ColName string
	Value   string
}

type ColData struct {
	ColName string
	ColType ColumnType
}

// Insert data to a table
func (table *Table) Insert(data []RowData) error {
	for _, col := range table.Columns {
		dataIndex := slices.IndexFunc(data, func(rowData RowData) bool { return rowData.ColName == col.Name })
		if dataIndex != -1 {
			value := data[dataIndex].Value
			col.Values = append(col.Values, value)
			continue
		}

		value, _ := col.Type.GetDefaultValue()
		col.Values = append(col.Values, value)
	}

	return nil
}

// Get data from a table
//   - columnNames define which columns to include, * get all.
//   - filters define which rows to include
func (table *Table) Get(columnNames []string, filters []*Filter) (*TableData, error) {
	rowCount := len(table.Columns[0].Values)
	columns := table.getColumns(columnNames)
	data := &TableData{
		Columns: Map(columns, func(col *Column) string { return col.Name }),
		Data:    [][]string{},
	}

	for rowIndex := 0; rowIndex < rowCount; rowIndex++ {
		if !table.isRowIncludedInFilters(rowIndex, filters) {
			continue
		}

		obj := make([]string, len(columns))
		for colIndex, col := range columns {
			value := col.Values[rowIndex]
			obj[colIndex] = value
		}

		data.Data = append(data.Data, obj)
	}

	return data, nil
}

// Update values of the table
func (table *Table) Update(data []RowData, filters []*Filter) error {
	colCount := len(table.Columns)
	rowCount := len(table.Columns[0].Values)

	for rowIndex := 0; rowIndex < rowCount; rowIndex++ {
		if !table.isRowIncludedInFilters(rowIndex, filters) {
			continue
		}

		for colIndex := 0; colIndex < colCount; colIndex++ {
			col := table.Columns[colIndex]
			dataIndex := slices.IndexFunc(data, func(valData RowData) bool { return valData.ColName == col.Name })

			if dataIndex != -1 {
				newValue := data[dataIndex].Value
				col.Values[rowIndex] = newValue
			}
		}
	}

	return nil
}

// Delete values from the table
func (table *Table) Delete(filters []*Filter) error {
	colCount := len(table.Columns)
	rowCount := len(table.Columns[0].Values)

	for rowIndex := rowCount - 1; rowIndex >= 0; rowIndex-- {
		if !table.isRowIncludedInFilters(rowIndex, filters) {
			continue
		}

		for colIndex := 0; colIndex < colCount; colIndex++ {
			col := table.Columns[colIndex]
			col.Values = slices.Delete(col.Values, rowIndex, rowIndex+1)
		}
	}

	return nil
}

// Get a column by name
func (table *Table) getColumnByName(colName string) (*Column, error) {
	index := slices.IndexFunc(table.Columns, func(col *Column) bool { return col.Name == colName })
	if index == -1 {
		return nil, fmt.Errorf("no column was found: %s", colName)
	}

	return table.Columns[index], nil
}

// Get all columns by name array, if name array contains only a single asterisk all columns are returned
func (table *Table) getColumns(columnNames []string) []*Column {
	if len(columnNames) == 1 && columnNames[0] == "*" {
		return table.Columns // SELECT ALL
	}

	columns := []*Column{}
	for _, colName := range columnNames {
		column, err := table.getColumnByName(colName)
		if err != nil {
			continue
		}

		columns = append(columns, column)
	}

	return columns
}

// Check if row is included in the filters
func (table *Table) isRowIncludedInFilters(rowIndex int, filters []*Filter) bool {
	for _, col := range table.Columns {
		value := col.Values[rowIndex]
		if IsTrueForAny(filters, func(filter *Filter) bool { return filter.ColumnName == col.Name && !filter.IsIncluded(value, col.Type) }) {
			// ROW NOT INCLUDED IN THE FILTERS
			return false
		}
	}

	return true
}
