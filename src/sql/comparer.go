package sql

import "strconv"

type EqualityOperator int

const (
	LESS             EqualityOperator = 1
	GREATER          EqualityOperator = 2
	EQUAL            EqualityOperator = 4
	LESS_OR_EQUAL    EqualityOperator = LESS | EQUAL
	GREATER_OR_EQUAL EqualityOperator = GREATER | EQUAL
)

func GetOperator(s string) EqualityOperator {
	switch s {
	case "<":
		return LESS
	case ">":
		return GREATER
	case "=":
		return EQUAL
	case "<=":
		return LESS_OR_EQUAL
	case ">=":
		return GREATER_OR_EQUAL
	}

	return -1
}

func (operator EqualityOperator) Compare(t ColumnType, a string, b string) bool {
	switch t {
	case TYPE_INT:
		intValue, _ := strconv.Atoi(a)
		intCompareValue, _ := strconv.Atoi(b)

		switch operator {
		case LESS:
			return intValue < intCompareValue
		case LESS_OR_EQUAL:
			return intValue <= intCompareValue
		case EQUAL:
			return intValue == intCompareValue
		case GREATER:
			return intValue > intCompareValue
		case GREATER_OR_EQUAL:
			return intValue >= intCompareValue
		}
	case TYPE_VARCHAR:
		if operator == EQUAL {
			return a == b
		}
	}

	return false
}
