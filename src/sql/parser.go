package sql

import (
	"fmt"
	"strings"
)

// Parse tokens to sql expression
func Parse(tokens []*Token) (Operation, error) {
	if len(tokens) <= 0 {
		return nil, fmt.Errorf("any operation could not be created, no tokens")
	}

	switch strings.ToUpper(tokens[0].Value) {
	case "SELECT":
		return parseSelect(tokens, 1)
	case "CREATE":
		return parseCreate(tokens, 1)
	case "INSERT":
		return parseInsert(tokens, 1)
	case "UPDATE":
		return parseUpdate(tokens, 1)
	case "DELETE":
		return parseDelete(tokens, 1)
	case "DROP":
		return parseDrop(tokens, 1)
	}

	return nil, fmt.Errorf("any operation could not be created, invalid or not supported operation")
}

// Parse a select operation
func parseSelect(tokens []*Token, index int) (Operation, error) {
	columnNames := []string{}

	for i := index; i < len(tokens); i += 2 {
		columnName := tokens[i].Value
		columnNames = append(columnNames, columnName)

		if tokens[i+1].Type == TOKEN_COMMA {
			continue
		}

		index = i + 1
		break
	}

	if strings.ToUpper(tokens[index].Value) != "FROM" {
		return nil, fmt.Errorf("select operation could not be created, missing from keyword")
	}

	tableName := tokens[index+1].Value
	index += 2

	filters := []*Filter{}
	sorters := []*Sorter{}
	for index < len(tokens) {
		switch strings.ToUpper(tokens[index].Value) {
		case "WHERE":
			f, i := parseFilter(tokens, index+1)
			filters = append(filters, f...)
			index = i
			continue
		case "ORDER":
			s, i, err := parseSorter(tokens, index+1)

			if err != nil {
				return nil, err
			}

			sorters = append(sorters, s)
			index = i
			continue
		}

		return nil, fmt.Errorf("select operation could not be created, invalid syntax after table name")
	}

	return &SelectOperation{
		ColumnNames: columnNames,
		TableName:   tableName,
		Filters:     filters,
		Sorters:     sorters,
	}, nil
}

// Parse a create operation
func parseCreate(tokens []*Token, index int) (Operation, error) {
	if strings.ToUpper(tokens[index].Value) != "TABLE" || tokens[index+2].Value != "(" {
		return nil, fmt.Errorf("create operation could not be created, trying to create something other than a table or missing parentheses")
	}

	tableName := tokens[index+1].Value
	data := []ColData{}

	for i := index + 3; i < len(tokens); i += 3 {
		colName := tokens[i].Value
		colType, _ := GetType(tokens[i+1].Value)
		data = append(data, ColData{ColName: colName, ColType: colType})

		if tokens[i+2].Type == TOKEN_COMMA {
			continue
		}

		if tokens[i+2].Value == ")" {
			break
		}
	}

	return &CreateOperation{
		TableName: tableName,
		Data:      data,
	}, nil
}

// Parse an insert operation
func parseInsert(tokens []*Token, index int) (Operation, error) {
	if strings.ToUpper(tokens[index].Value) != "INTO" || tokens[index+2].Value != "(" {
		return nil, fmt.Errorf("insert operation could not be created, missing into keyword or parentheses")
	}

	tableName := tokens[index+1].Value
	data := []RowData{}

	for i := index + 3; i < len(tokens); i += 2 {
		columnName := tokens[i].Value
		data = append(data, RowData{ColName: columnName})

		if tokens[i+1].Type == TOKEN_COMMA {
			continue
		}

		if tokens[i+1].Value == ")" {
			index = i + 2
			break
		}
	}

	if strings.ToUpper(tokens[index].Value) != "VALUES" || tokens[index+1].Value != "(" {
		return nil, fmt.Errorf("insert operation could not be created, missing values keyword or parenthesis")
	}

	valIndex := 0
	for i := index + 2; i < len(tokens); i += 2 {
		value := tokens[i].Value
		data[valIndex].Value = value

		valIndex++

		if tokens[i+1].Type == TOKEN_COMMA {
			continue
		}

		if tokens[i+1].Value == ")" {
			break
		}
	}

	return &InsertOperation{
		TableName: tableName,
		Data:      data,
	}, nil
}

// Parse an update operation
func parseUpdate(tokens []*Token, index int) (Operation, error) {
	if tokens[index+1].Value != "(" {
		return nil, fmt.Errorf("update operation could not be created, missing parenthesis")
	}

	tableName := tokens[index].Value
	data := []RowData{}

	for i := index + 2; i < len(tokens); i += 2 {
		columnName := tokens[i].Value
		data = append(data, RowData{ColName: columnName})

		if tokens[i+1].Type == TOKEN_COMMA {
			continue
		}

		if tokens[i+1].Value == ")" {
			index = i + 2
			break
		}
	}

	if strings.ToUpper(tokens[index].Value) != "VALUES" || tokens[index+1].Value != "(" {
		return nil, fmt.Errorf("update operation could not be created, missing values keyword or parenthesis")
	}

	valIndex := 0
	for i := index + 2; i < len(tokens); i += 2 {
		value := tokens[i].Value
		data[valIndex].Value = value
		valIndex++

		if tokens[i+1].Type == TOKEN_COMMA {
			continue
		}

		if tokens[i+1].Value == ")" {
			index = i + 2
			break
		}
	}

	filters := []*Filter{}
	for index < len(tokens) {
		switch strings.ToUpper(tokens[index].Value) {
		case "WHERE":
			f, i := parseFilter(tokens, index+1)
			filters = append(filters, f...)
			index = i
			continue
		}

		return nil, fmt.Errorf("update operation could not be created, invalid syntax after table name")
	}

	return &UpdateOperation{
		TableName: tableName,
		Data:      data,
		Filters:   filters,
	}, nil
}

// Parse a delete operation
func parseDelete(tokens []*Token, index int) (Operation, error) {
	if strings.ToUpper(tokens[index].Value) != "FROM" {
		return nil, fmt.Errorf("delete operation could not be created, missing from keyword")
	}

	tableName := tokens[index+1].Value
	filters := []*Filter{}
	index += 2

	for index < len(tokens) {
		switch strings.ToUpper(tokens[index].Value) {
		case "WHERE":
			f, i := parseFilter(tokens, index+1)
			filters = append(filters, f...)
			index = i
			continue
		}

		return nil, fmt.Errorf("delete operation could not be created, invalid syntax after table name")
	}

	return &DeleteOperation{
		TableName: tableName,
		Filters:   filters,
	}, nil
}

// Parse a drop operation
func parseDrop(tokens []*Token, index int) (Operation, error) {
	if strings.ToUpper(tokens[index].Value) != "TABLE" {
		return nil, fmt.Errorf("drop operation could not be created, missing table keyword")
	}

	tableName := tokens[index+1].Value
	return &DropOperation{
		TableName: tableName,
	}, nil
}

// Parse a where expression.
// for example 0 < x < 1 returns two filters x > 0 and x < 1
func parseFilter(tokens []*Token, index int) ([]*Filter, int) {
	filters := []*Filter{}
	value1 := tokens[index].Value

	index++
	operator1 := GetEqualityOperator(tokens[index].Value)
	if tokens[index+1].Type == TOKEN_OPERATOR {
		operator1 = operator1 | GetEqualityOperator(tokens[index+1].Value)
		index++
	}

	index++
	value2 := tokens[index].Value

	index++
	if index < len(tokens) && tokens[index].Type == TOKEN_OPERATOR {
		operator2 := GetEqualityOperator(tokens[index].Value)
		if tokens[index+1].Type == TOKEN_OPERATOR {
			operator2 = operator2 | GetEqualityOperator(tokens[index+1].Value)
			index++
		}

		index++
		value3 := tokens[index].Value

		filters = append(filters, &Filter{
			ColumnName:   value2,
			Operator:     operator1.Inverse(),
			CompareValue: value1,
		}, &Filter{
			ColumnName:   value2,
			Operator:     operator2,
			CompareValue: value3,
		})

		return filters, index + 1
	}

	filters = append(filters, &Filter{
		ColumnName:   value1,
		Operator:     operator1,
		CompareValue: value2,
	})

	return filters, index + 1
}

// Parse order by expression
func parseSorter(tokens []*Token, index int) (*Sorter, int, error) {
	if strings.ToUpper(tokens[index].Value) != "BY" {
		return nil, -1, fmt.Errorf("order could not be created, missing by keyword")
	}

	columnName := tokens[index+1].Value
	direction, err := GetSortDirection(tokens[index+2].Value)
	if err != nil {
		return nil, -1, err
	}

	return &Sorter{
		ColumnName: columnName,
		Direction:  direction,
	}, index + 3, nil
}
