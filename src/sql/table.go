package sql

import (
	"slices"
)

type Table struct {
	Name    string    `json:"table"`
	Columns []*Column `json:"columns"`
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
			column := table.Columns[col]
			value := column.Values[row]

			if !matchFilters(value, column, filters) {
				// OBJECT NOT INCLUDED IN THE FILTERS
				isIncluded = false
				break
			}

			if matchSelection(column, columnNames) {
				// COLUMN INCLUDED IN THE SELECTION
				obj.Values = append(obj.Values, value)
			}
		}

		if isIncluded {
			objects = append(objects, obj)
		}
	}

	return objects
}

func (table *Table) Update(data []ValData, filters []*Filter) {
	colCount := len(table.Columns)
	rowCount := len(table.Columns[0].Values)

	for row := 0; row < rowCount; row++ {
		isIncluded := true

		if len(filters) > 0 {
			for col := 0; col < colCount; col++ {
				column := table.Columns[col]
				value := column.Values[row]

				if !matchFilters(value, column, filters) {
					isIncluded = false
					break
				}
			}
		}

		if !isIncluded {
			// OBJECT NOT INCLUDED IN THE FILTERS
			continue
		}

		for col := 0; col < colCount; col++ {
			column := table.Columns[col]
			valIndex := slices.IndexFunc(data, func(valData ValData) bool { return valData.ColumnName == column.Name })

			if valIndex != -1 {
				value := data[valIndex].Value
				column.Values[row] = value
			}
		}
	}
}

func (table *Table) Delete(filters []*Filter) {
	colCount := len(table.Columns)
	rowCount := len(table.Columns[0].Values)

	for row := rowCount - 1; row >= 0; row-- {
		isIncluded := true

		if len(filters) > 0 {
			for col := 0; col < colCount; col++ {
				column := table.Columns[col]
				value := column.Values[row]

				if !matchFilters(value, column, filters) {
					isIncluded = false
					break
				}
			}
		}

		if !isIncluded {
			// OBJECT NOT INCLUDED IN THE FILTERS
			continue
		}

		for col := 0; col < colCount; col++ {
			column := table.Columns[col]
			column.Values = slices.Delete(column.Values, row, row+1)
		}
	}
}

func matchSelection(column *Column, columnNames []string) bool {
	if len(columnNames) == 1 && columnNames[0] == "*" {
		return true // SELECT ALL
	}

	return slices.ContainsFunc(columnNames, func(columnName string) bool {
		return columnName == column.Name
	})
}

func matchFilters(value string, column *Column, filters []*Filter) bool {
	return !IsTrueForAny(filters, func(filter *Filter) bool {
		return filter.ColumnName == column.Name && !filter.IsIncluded(value, column.Type)
	})
}
