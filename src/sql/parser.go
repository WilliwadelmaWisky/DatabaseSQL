package sql

import (
	"fmt"
	"strings"
)

type ColData struct {
	Name string
	Type ColumnType
}

type ValData struct {
	ColumnName string
	Value      string
}

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
	}

	return nil, fmt.Errorf("any operation could not be created, invalid or not supported operation")
}

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
			f, i := parseFilters(tokens, index+1)
			filters = append(filters, f...)
			index = i
			continue
		case "ORDER":
			s, i := parseOrder(tokens, index+1)
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

func parseCreate(tokens []*Token, index int) (Operation, error) {
	if strings.ToUpper(tokens[index].Value) != "TABLE" || tokens[index+2].Value != "(" {
		return nil, fmt.Errorf("create operation could not be created, trying to create something other than a table or missing parentheses")
	}

	tableName := tokens[index+1].Value
	data := []ColData{}

	for i := index + 3; i < len(tokens); i += 3 {
		colName := tokens[i].Value
		colType, _ := GetType(tokens[i+1].Value)
		data = append(data, ColData{Name: colName, Type: colType})

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

func parseInsert(tokens []*Token, index int) (Operation, error) {
	if strings.ToUpper(tokens[index].Value) != "INTO" || tokens[index+2].Value != "(" {
		return nil, fmt.Errorf("insert operation could not be created, missing into keyword or parentheses")
	}

	tableName := tokens[index+1].Value
	data := []ValData{}

	for i := index + 3; i < len(tokens); i += 2 {
		columnName := tokens[i].Value
		data = append(data, ValData{ColumnName: columnName})

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

func parseUpdate(tokens []*Token, index int) (Operation, error) {
	if tokens[index+1].Value != "(" {
		return nil, fmt.Errorf("update operation could not be created, missing parenthesis")
	}

	tableName := tokens[index].Value
	data := []ValData{}

	for i := index + 2; i < len(tokens); i += 2 {
		columnName := tokens[i].Value
		data = append(data, ValData{ColumnName: columnName})

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
			f, i := parseFilters(tokens, index+1)
			filters = append(filters, f...)
			index = i
			continue
		}

		return nil, fmt.Errorf("select operation could not be created, invalid syntax after table name")
	}

	return &UpdateOperation{
		TableName: tableName,
		Data:      data,
		Filters:   filters,
	}, nil
}

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
			f, i := parseFilters(tokens, index+1)
			filters = append(filters, f...)
			index = i
			continue
		}

		return nil, fmt.Errorf("select operation could not be created, invalid syntax after table name")
	}

	return &DeleteOperation{
		TableName: tableName,
		Filters:   filters,
	}, nil
}

func parseFilters(tokens []*Token, index int) ([]*Filter, int) {
	filters := []*Filter{}
	value1 := tokens[index].Value

	index++
	operator1 := GetOperator(tokens[index].Value)
	if tokens[index+1].Type == TOKEN_OPERATOR {
		operator1 = operator1 | GetOperator(tokens[index+1].Value)
		index++
	}

	index++
	value2 := tokens[index].Value

	filters = append(filters, &Filter{
		ColumnName:   value1,
		Operator:     operator1,
		CompareValue: value2,
	})

	// index++
	// if index < len(tokens) && tokens[index].Type == TOKEN_OPERATOR {
	// 	operator2 := stringToEqualityOperator(tokens[index].Value)
	// 	if tokens[index+1].Type == TOKEN_OPERATOR {
	// 		operator2 = operator2 | stringToEqualityOperator(tokens[index+1].Value)
	// 		index++
	// 	}

	// 	index++
	// 	value3 := tokens[index].Value

	// }

	index++
	return filters, index
}

func parseOrder(tokens []*Token, index int) (*Sorter, int) {
	return &Sorter{}, 0
}
