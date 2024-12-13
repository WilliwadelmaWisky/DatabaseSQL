package sql

import (
	"encoding/json"
	"fmt"
	"slices"
)

type Create struct {
	Table   string
	Columns []Column
}

type Insert struct {
	Table     string
	Variables []Variable
}

type Update struct {
	Table     string
	Variables []Variable
}

type Delete struct {
	Table string
}

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

type Type int

const (
	INT Type = iota
	VARCHAR
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

type Operation interface {
	Call(database *Database) []byte
}

type SelectOperation struct {
	Columns []string
	Table   string
}

func (operation *SelectOperation) Call(database *Database) []byte {
	table, _ := database.GetTable(operation.Table)
	values := table.Get(func(column *Column) bool {
		if len(operation.Columns) == 1 && operation.Columns[0] == "*" {
			return true // SELECT ALL
		}

		return slices.ContainsFunc(operation.Columns, func(columnName string) bool { return columnName == column.Name })
	})

	bytes, _ := json.Marshal(values)
	return bytes
}

type CreateOperation struct {
	Table string
}

func (operation *CreateOperation) Call(database *Database) []byte {
	return []byte{}
}

type InsertOperation struct {
	Table string
}

func (operation *InsertOperation) Call(database *Database) []byte {
	return []byte{}
}

func Parse(tokens []Token) (Operation, error) {
	if len(tokens) <= 0 {
		return nil, fmt.Errorf("no operation could be created, no tokens!")
	}

	if tokens[0].Value == "SELECT" {
		col := []string{}
		table := ""

		for i := 1; i < len(tokens); i++ {
			if tokens[i].Type != Comma {
				col = append(col, tokens[i].Value)
			}

			if tokens[i].Value == "FROM" && i < len(tokens) {
				i++
				table = tokens[i].Value
			}
		}

		return &SelectOperation{
			Columns: col,
			Table:   table,
		}, nil
	}

	if tokens[0].Value == "CREATE" && tokens[1].Value == "TABLE" && tokens[3].Type == Parenthesis {
		table := tokens[2].Value

		columns := []string{}
		for i := 4; i < len(tokens); i += 3 {
			column := tokens[i].Value
			colType := tokens[i+1].Value
			columns = append(columns, column)

			if tokens[i+2].Type == Parenthesis {
				break
			}

			if tokens[i+2].Type == Comma {
				continue
			}
		}

		return &CreateOperation{
			Table: table,
		}, nil
	}

	if tokens[0].Value == "INSERT" && tokens[1].Value == "INTO" && tokens[4].Type == Parenthesis {
		table := tokens[2].Value
		data := []struct {
			Col string
			Val string
		}{}

		index := 4
		for i := index; i < len(tokens); i += 2 {
			column := tokens[i].Value
			data = append(data, struct {
				Col string
				Val string
			}{Col: column})

			if tokens[i+1].Type == Comma {
				continue
			}

			if tokens[i+1].Type == Parenthesis {
				if tokens[i+2].Value == "VALUES" && tokens[i+3].Type == Parenthesis {
					index = i + 4
				}

				break
			}
		}

		valIndex := 0
		for i := index; i < len(tokens); i += 2 {
			value := tokens[i].Value
			data[valIndex].Val = value

			valIndex++

			if tokens[i+1].Type == Comma {
				continue
			}

			if tokens[i+1].Type == Parenthesis {
				break
			}
		}

		return &InsertOperation{
			Table: table,
		}, nil
	}

	return nil, fmt.Errorf("no operation could be created, invalid or not supported tokens!")
}
