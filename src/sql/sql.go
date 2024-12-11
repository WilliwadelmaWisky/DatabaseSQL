package sql

type Select struct {
	Columns []string
	Table   string
}

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

type Column struct {
	Name   string
	Type   Type
	Length int
}

type Variable struct {
	Column string
	Value  string
}

type Operation int

const (
	CREATE Operation = iota
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

func Parse(tokens []Token) {
	if len(tokens) <= 0 {
		return
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

		sel := Select{
			Columns: col,
			Table:   table,
		}
	}
}
