package sql

import (
	"slices"
)

type Table struct {
	Name    string    `json:"table"`
	Columns []*Column `json:"columns"`
}

func (table *Table) Get(columnNames []string, filters []*Filter) []*Object {
	colCount := len(table.Columns)
	rowCount := len(table.Columns[0].Values)
	objects := []*Object{}

	// TABLE DEBUGGING
	// fmt.Printf("(%d x %d) -matrix\n", rowCount, colCount)

	for row := 0; row < rowCount; row++ {
		obj := &Object{}
		isIncluded := true

		for col := 0; col < colCount; col++ {
			value := table.Columns[col].Values[row]

			if IsTrueForAny(filters, func(filter *Filter) bool {
				return filters[0].ColumnName == table.Columns[col].Name && !filters[0].IsIncluded(value, table.Columns[col].Type)
			}) {
				// OBJECT NOT INCLUDED IN THE FILTERS
				isIncluded = false
				break
			}

			if isColumnIncluded(table.Columns[col], columnNames) {
				obj.Values = append(obj.Values, value)
			}
		}

		if isIncluded {
			objects = append(objects, obj)
		}
	}

	return objects
}

func (table *Table) Insert(data []ValData) {
	for _, column := range table.Columns {
		valIndex := slices.IndexFunc(data, func(valData ValData) bool { return valData.ColumnName == column.Name })
		if valIndex != -1 {
			value := data[valIndex].Value
			column.Values = append(column.Values, value)
			continue
		}

		value, _ := column.Type.GetDefaultValue()
		column.Values = append(column.Values, value)
	}
}

func isColumnIncluded(column *Column, columnNames []string) bool {
	if len(columnNames) == 1 && columnNames[0] == "*" {
		return true // SELECT ALL
	}

	return slices.ContainsFunc(columnNames, func(columnName string) bool { return columnName == column.Name })
}
