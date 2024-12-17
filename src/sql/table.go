package sql

import (
	"fmt"
	"slices"
)

type Table struct {
	Name    string    `json:"table"`
	Columns []*Column `json:"columns"`
}

type TableData struct {
	ColumnNames []string   `json:"column_names"`
	Rows        [][]string `json:"rows"`
	RowCount    int        `json:"row_count"`
}

type RowData struct {
	ColName string
	Value   string
}

type ColData struct {
	ColName string
	ColType ColumnType
}

func (table *Table) Insert(data []RowData) {
	for _, column := range table.Columns {
		valIndex := slices.IndexFunc(data, func(valData RowData) bool { return valData.ColName == column.Name })
		if valIndex != -1 {
			value := data[valIndex].Value
			column.Values = append(column.Values, value)
			continue
		}

		value, _ := column.Type.GetDefaultValue()
		column.Values = append(column.Values, value)
	}
}

func (table *Table) Get(columnNames []string, filters []*Filter) *TableData {
	rowCount := len(table.Columns[0].Values)
	columns := table.getColumns(columnNames)
	data := &TableData{
		ColumnNames: Map(columns, func(col *Column) string { return col.Name }),
		Rows:        [][]string{},
		RowCount:    0,
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

		data.Rows = append(data.Rows, obj)
		data.RowCount++
	}

	return data
}

func (table *Table) Update(data []RowData, filters []*Filter) {
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
}

func (table *Table) Delete(filters []*Filter) {
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
}

func (table *Table) getColumnByName(colName string) (*Column, error) {
	index := slices.IndexFunc(table.Columns, func(col *Column) bool { return col.Name == colName })
	if index == -1 {
		return nil, fmt.Errorf("no column was found: %s", colName)
	}

	return table.Columns[index], nil
}

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
