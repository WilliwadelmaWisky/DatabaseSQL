package sql

import (
	"fmt"
)

type Variable struct {
	Column string
	Value  string
}

type OperationType int

const (
	CREATE OperationType = iota
	INSERT
	UPDATE
	SELECT
	DELETE
)

type EqualityOperator byte

const (
	LESS             EqualityOperator = 1
	GREATER          EqualityOperator = 2
	EQUAL            EqualityOperator = 4
	LESS_OR_EQUAL    EqualityOperator = LESS | EQUAL
	GREATER_OR_EQUAL EqualityOperator = GREATER | EQUAL
)

type Where struct {
	Column    string
	Condition []Condition
}

type Condition struct {
	Operator EqualityOperator
}

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

func Parse(tokens []Token) (Operation, error) {
	if len(tokens) <= 0 {
		return nil, fmt.Errorf("no operation could be created, no tokens")
	}

	if tokens[0].Value == "SELECT" {
		columnNames := []string{}
		index := 1

		for i := index; i < len(tokens); i += 2 {
			columnName := tokens[i].Value
			columnNames = append(columnNames, columnName)

			if tokens[i+1].Type == Comma {
				continue
			}

			index = i + 1
			break
		}

		tableName := ""
		if tokens[index].Value == "FROM" {
			tableName = tokens[index+1].Value
		}

		return &SelectOperation{
			ColumnNames: columnNames,
			TableName:   tableName,
		}, nil
	}

	if tokens[0].Value == "CREATE" && tokens[1].Value == "TABLE" && tokens[3].Value == "(" {
		tableName := tokens[2].Value
		data := []ColData{}

		for i := 4; i < len(tokens); i += 3 {
			colName := tokens[i].Value
			colType, _ := getColumnType(tokens[i+1].Value)
			data = append(data, ColData{Name: colName, Type: colType})

			if tokens[i+2].Value == ")" {
				break
			}

			if tokens[i+2].Type == Comma {
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

			if tokens[i+1].Type == Comma {
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

			if tokens[i+1].Type == Comma {
				continue
			}

			if tokens[i+1].Type == Parenthesis {
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

func getColumnType(typeString string) (ColumnType, error) {
	switch typeString {
	case "INT":
		return TYPE_INT, nil
	case "VARCHAR":
		return TYPE_VARCHAR, nil
	}

	return TYPE_INVALID, fmt.Errorf("invalid column type [%s]", typeString)
}
