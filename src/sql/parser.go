package sql

import (
	"fmt"
)

type OrderDirection int

const (
	ASC OrderDirection = iota
	DESC
)

type Order struct {
	Column    string
	Direction OrderDirection
}

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
		return nil, fmt.Errorf("no operation could be created, no tokens")
	}

	if tokens[0].Value == "SELECT" {
		columnNames := []string{}
		index := 1

		for i := index; i < len(tokens); i += 2 {
			columnName := tokens[i].Value
			columnNames = append(columnNames, columnName)

			if tokens[i+1].Type == TOKEN_COMMA {
				continue
			}

			index = i + 1
			break
		}

		tableName := ""
		if tokens[index].Value == "FROM" {
			tableName = tokens[index+1].Value
			index += 2
		}

		filters, indexIncrement := parseFilters(tokens[index:])
		index += indexIncrement

		return &SelectOperation{
			ColumnNames: columnNames,
			TableName:   tableName,
			Filters:     filters,
		}, nil
	}

	if tokens[0].Value == "CREATE" && tokens[1].Value == "TABLE" && tokens[3].Value == "(" {
		tableName := tokens[2].Value
		data := []ColData{}

		for i := 4; i < len(tokens); i += 3 {
			colName := tokens[i].Value
			colType, _ := GetType(tokens[i+1].Value)
			data = append(data, ColData{Name: colName, Type: colType})

			if tokens[i+2].Value == ")" {
				break
			}

			if tokens[i+2].Type == TOKEN_COMMA {
				continue
			}
		}

		return &CreateOperation{
			TableName: tableName,
			Data:      data,
		}, nil
	}

	if tokens[0].Value == "INSERT" && tokens[1].Value == "INTO" && tokens[3].Value == "(" {
		tableName := tokens[2].Value
		data := []ValData{}

		index := 4
		for i := index; i < len(tokens); i += 2 {
			columnName := tokens[i].Value
			data = append(data, ValData{ColumnName: columnName})

			if tokens[i+1].Type == TOKEN_COMMA {
				continue
			}

			if tokens[i+1].Value == ")" {
				if tokens[i+2].Value == "VALUES" && tokens[i+3].Value == "(" {
					index = i + 4
				}

				break
			}
		}

		valIndex := 0
		for i := index; i < len(tokens); i += 2 {
			value := tokens[i].Value
			data[valIndex].Value = value

			valIndex++

			if tokens[i+1].Type == TOKEN_COMMA {
				continue
			}

			if tokens[i+1].Type == TOKEN_PARENTHESIS {
				break
			}
		}

		return &InsertOperation{
			TableName: tableName,
			Data:      data,
		}, nil
	}

	return nil, fmt.Errorf("no operation could be created, invalid or not supported tokens")
}

func parseFilters(tokens []*Token) ([]*Filter, int) {
	filters := []*Filter{}
	if len(tokens) <= 0 || tokens[0].Value != "WHERE" {
		return filters, 0
	}

	index := 1
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

// func parseOrder() {

// }
