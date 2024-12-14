package sql

import (
	"slices"
)

type Table struct {
	Name    string    `json:"table"`
	Columns []*Column `json:"columns"`
}

func (table *Table) Get(match func(*Column) bool) []*Object {
	colCount := len(table.Columns)
	rowCount := len(table.Columns[0].Values)
	objects := make([]*Object, rowCount)

	// TABLE DEBUGGING
	// fmt.Printf("(%d x %d) -matrix\n", rowCount, colCount)

	for row := 0; row < rowCount; row++ {
		obj := &Object{
			Values: make([]string, colCount),
		}

		for col := 0; col < colCount; col++ {
			obj.Values[col] = table.Columns[col].Values[row]
		}

		objects[row] = obj
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
